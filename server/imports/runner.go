package imports

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/progress"
	"github.com/lsdch/biome/models"
	csvmodels "github.com/lsdch/biome/models/csv"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/services/storage"
	"github.com/lsdch/biome/stores"
	"github.com/sirupsen/logrus"
)

const (
	ErrGBIFAlreadyRunning string = "GBIF import is already in progress for this import hash"
)

type RunnerStatus string

const (
	Created            RunnerStatus = "created"
	Staging            RunnerStatus = "staging"
	Staged             RunnerStatus = "staged"
	Running            RunnerStatus = "running"
	NeedsResolution    RunnerStatus = "needs_resolution"
	ReadyToMaterialize RunnerStatus = "ready_to_materialize"
	Materializing      RunnerStatus = "materializing"
	Completed          RunnerStatus = "completed"
	Failed             RunnerStatus = "failed"
	Cancelled          RunnerStatus = "cancelled"
)

type ImportRunner struct {
	db            *db.DB
	batch         models.ImportBatch
	store         *stores.BatchesStore
	taxonResolver TaxonResolver
	bibliography  *BibliographyResolver
	samplings     *services.SamplingService
	occurrences   *services.OccurrencesService
	fileStorage   storage.RawFileStorage

	parser    CSVParser
	validator *validator.Validate

	status           RunnerStatus
	resolutionStatus models.MaterializationReadyCheck
	gbif             *progress.ProgressTracker

	mu         sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	err        error
	events     ImportEventSink[BatchSnapshot]
	lastNotify atomic.Value //time.Time
}

func NewImportRunner(ctx context.Context, db *db.DB,
	events ImportEventSink[BatchSnapshot],
	batch models.ImportBatch,
	store *stores.BatchesStore,
	samplings *services.SamplingService,
	taxonResolver TaxonResolver,
	bibliography *BibliographyResolver,
	occurrences *services.OccurrencesService,
	fileStorage storage.RawFileStorage,
) *ImportRunner {
	ctx, cancel := context.WithCancel(ctx)

	runner := &ImportRunner{
		db:            db,
		batch:         batch,
		store:         store,
		taxonResolver: taxonResolver,
		samplings:     samplings,
		bibliography:  bibliography,
		occurrences:   occurrences,
		status:        Created,
		parser:        NewCSVParser(),
		ctx:           ctx,
		cancel:        cancel,
		events:        events,
		validator:     validator.New(validator.WithRequiredStructEnabled()),
	}
	runner.gbif = progress.NewProgressTracker().WithCallback(func(status progress.ProgressStatus) {
		if status == progress.Running && time.Since(runner.lastNotificationTime()) < 1000*time.Millisecond {
			return
		}

		runner.notify()
	})
	return runner
}

func (r *ImportRunner) lastNotificationTime() time.Time {
	if t, ok := r.lastNotify.Load().(time.Time); ok {
		return t
	}
	return time.Time{}
}

func (r *ImportRunner) ID() uuid.UUID {
	return r.batch.ID
}

func (r *ImportRunner) Batch() models.ImportBatch {
	return r.batch
}

func (r *ImportRunner) Status() RunnerStatus {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.status
}

func (r *ImportRunner) StartBatchCSV(csvFile io.Reader, separator rune, taxonDefinitions []models.TaxonDefinition, mergeUndated bool) error {

	if err := r.StartStaging(); err != nil {
		return err
	}
	rows, err := r.parser.ParseCSV(csvFile, separator)
	if err != nil {
		return err
	}

	if err := csvmodels.ValidateRows(rows, r.validator); err != nil {
		r.Fail(err)
		return err
	}

	var bib = []csvmodels.PublicationImportRow{}
	for _, row := range rows {
		if row.HasPublication() {
			pub := row.Publication
			pub.SetRowNumber(row.RowNumber())
			bib = append(bib, pub)
		}
	}

	// bib, err := r.parser.ParseBibCSV(csvFile, separator)
	// if err != nil {
	// 	return err
	// }

	if err := csvmodels.ValidateRows(bib, r.validator); err != nil {
		r.Fail(err)
		return err
	}

	err = r.db.WithTx(r.ctx, func(tx *db.Tx) error {
		if err := r.store.InsertStaging(r.ctx, tx, r.batch.ID, rows, taxonDefinitions, mergeUndated); err != nil {
			return fmt.Errorf("insert staging occurrences: %w", err)
		}
		if _, err = r.taxonResolver.InitResolution(r.ctx, tx, r.batch.ID); err != nil {
			return fmt.Errorf("init taxon resolution: %w", err)
		}
		if _, err := r.samplings.InitMethodResolution(r.ctx, tx, r.batch.ID); err != nil {
			return fmt.Errorf("init method resolution: %w", err)
		}
		if _, err := r.samplings.InitFixativeResolution(r.ctx, tx, r.batch.ID); err != nil {
			return fmt.Errorf("init fixative resolution: %w", err)
		}

		if err := r.bibliography.InitBibliographyResolution(r.ctx, tx, r.batch.ID, bib); err != nil {
			return fmt.Errorf("init bibliography resolution: %w", err)
		}

		if err := r.MarkStaged(tx); err != nil {
			return fmt.Errorf("mark staged: %w", err)
		}
		return nil
	})
	if err != nil {
		r.Fail(err)
		return err
	}

	go r.Run()

	return nil
}

func (r *ImportRunner) AddBibliographyCSV(csvFile io.Reader, separator rune) error {
	rows, err := r.parser.ParseBibCSV(csvFile, separator)
	if err != nil {
		return err
	}

	if err := csvmodels.ValidateRows(rows, r.validator); err != nil {
		r.Fail(err)
		return err
	}

	err = r.db.WithTx(r.ctx, func(tx *db.Tx) error {
		if err := r.bibliography.InitBibliographyResolution(r.ctx, tx, r.batch.ID, rows); err != nil {
			return fmt.Errorf("init bibliography resolution: %w", err)
		}
		return nil
	})
	if err != nil {
		r.Fail(err)
		return err
	}

	go r.Run()

	return nil
}

func (r *ImportRunner) Runnable() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.status != Running && r.status != Materializing && r.status != Completed
	// return r.status == Staged || r.status == NeedsResolution || r.status == Failed || r.status == Cancelled
	// return true
}

func (r *ImportRunner) Run() error {
	wg := sync.WaitGroup{}
	if r.Runnable() {
		r.withMutex(func() {
			r.status = Running
			r.err = nil
		})
		if r.gbif.Snapshot().Status != progress.Running {
			wg.Add(1)
			go func() {
				defer wg.Done()
				r.EnrichGBIF()
			}()
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.EnrichBibliography()

			logrus.Infof("Running automatic bibliography resolution for unambiguous candidates for import ID %s", r.batch.ID)
			if err := r.bibliography.AutoResolveBibliography(r.ctx, r.db, r.batch.ID, 80, 10); err != nil {
				logrus.Errorf("autoresolve bibliography failed for import ID %s: %v", r.batch.ID, err)
				r.Fail(fmt.Errorf("autoresolve bibliography failed: %w", err))
				return
			}

			if err := r.bibliography.ResolveBibliographyManualCandidates(r.ctx, r.db, r.batch.ID); err != nil {
				logrus.Errorf("manual resolve bibliography failed for import ID %s: %v", r.batch.ID, err)
				r.Fail(fmt.Errorf("manual resolve bibliography failed: %w", err))
				return
			}
		}()
	}
	wg.Wait()

	r.CheckReadyToMaterialize()
	return nil
}

func (r *ImportRunner) EnrichBibliography() {

	toStage, err := r.bibliography.FetchExternalCandidatesDOI(r.ctx, r.db, r.batch.ID)
	if err != nil {
		logrus.Errorf("error fetching external candidates for bibliography import %s: %v", r.batch.ID, err)
		r.Fail(err)
	}
	// verbatimToStage, err := r.bibliography.FetchExternalCandidatesQuery(r.ctx, r.db, r.batch.ID)
	// if err != nil {
	// 	logrus.Errorf("error fetching external candidates for bibliography import %s: %v", r.batch.ID, err)
	// 	r.Fail(err)
	// }

	// toStage := append(doisToStage, verbatimToStage...)

	logrus.Debugf("Staging %d external candidates for bibliography import %s", len(toStage), r.batch.ID)
	if err := r.bibliography.StageExternalCandidates(r.ctx, r.db, r.batch.ID, toStage); err != nil {
		logrus.Errorf("error staging external candidates for bibliography import %s: %v", r.batch.ID, err)
		r.Fail(err)
	}
}

func (r *ImportRunner) enrichGBIF() error {
	toFetch, err := r.taxonResolver.ListTaxaToFetch(r.ctx, r.db, r.batch.ID)
	if err != nil {
		return err
	} else if len(toFetch) == 0 {
		// Enrichment is complete when there are no more taxa to fetch from GBIF
		r.gbif.Complete()
		return nil
	}

	logrus.Infof("Fetching %d taxa from GBIF for import ID %s", len(toFetch), r.batch.ID)

	taxa, err := r.taxonResolver.FetchCandidatesFromGBIF(r.ctx, r.batch.TaxonomicScope, toFetch, r.gbif)
	if err != nil {
		return err
	}

	err = r.taxonResolver.InsertGBIFCandidates(r.ctx, r.db, r.batch.ID, taxa)
	if err != nil {
		return err
	}

	err = r.taxonResolver.MarkTaxaGBIFImportCompleted(r.ctx, r.db, r.batch.ID)
	if err != nil {
		return err
	}

	logrus.Infof("Running automatic taxon resolution for unambiguous candidates for import ID %s", r.batch.ID)
	if err := r.taxonResolver.AutoResolveUnambiguousCandidates(r.ctx, r.db, r.batch.ID); err != nil {
		logrus.Errorf("autoresolve failed for import ID %s: %v", r.batch.ID, err)
		return fmt.Errorf("autoresolve failed: %w", err)
	}

	logrus.Infof("Running automatic creation of manual candidates for import ID %s", r.batch.ID)
	if err := r.taxonResolver.AutoCreateManualCandidates(r.ctx, r.db, r.batch.ID); err != nil {
		logrus.Errorf("auto-create manual candidates failed for import ID %s: %v", r.batch.ID, err)
		return fmt.Errorf("auto-create manual candidates failed: %w", err)
	}

	logrus.Infof("Resolving sampling targets for import ID %s", r.batch.ID)
	if err := r.taxonResolver.InitSamplingTargetsResolution(r.ctx, r.db, r.batch.ID); err != nil {
		logrus.Errorf("init sampling targets resolution failed for import ID %s: %v", r.batch.ID, err)
		return fmt.Errorf("init sampling targets resolution failed: %w", err)
	}

	logrus.Infof("Running automatic taxon resolution for unambiguous candidates after manual candidate creation for import ID %s", r.batch.ID)
	if err := r.taxonResolver.AutoResolveUnambiguousCandidates(r.ctx, r.db, r.batch.ID); err != nil {
		logrus.Errorf("autoresolve failed for import ID %s: %v", r.batch.ID, err)
		return fmt.Errorf("autoresolve failed: %w", err)
	}
	// Restart the GBIF enrichment process to fetch new candidates for any new resolutions *
	// created by the auto-create manual candidates step
	return r.enrichGBIF()
}

func (r *ImportRunner) EnrichGBIF() {

	if r.gbif.Snapshot().Status != progress.Running {
		r.gbif.Start(0)
	}

	if err := r.enrichGBIF(); err != nil {
		r.Fail(err)
		r.gbif.Fail(err)
		return
	}

	if err := r.taxonResolver.SetNeedsResolution(r.ctx, r.db, r.batch.ID); err != nil {
		r.Fail(err)
		r.gbif.Fail(err)
		return
	}
	r.gbif.Complete()
}

func (r *ImportRunner) withMutex(fn func()) {
	r.mu.Lock()
	fn()
	r.mu.Unlock()
	r.notify()
}

func (r *ImportRunner) StartStaging() error {
	if r.Status() == Created {
		r.withMutex(func() {
			r.status = Staging
			r.err = nil
		})
	} else {
		return fmt.Errorf("cannot start batch in status %s", r.status)
	}
	return nil
}

func (r *ImportRunner) MarkStaged(db db.Querier) error {
	if err := r.store.SetBatchStatus(r.ctx, db, r.batch.ID, biomedb.ImportBatchStatusStaged); err != nil {
		r.Fail(err)
		return err
	}
	r.withMutex(func() {
		r.status = Staged
	})
	return nil
}

func (r *ImportRunner) Fail(err error) {

	logrus.Errorf("Import runner failed for import ID %s: %v", r.batch.ID, err)

	if errors.Is(err, context.Canceled) {
		if err := r.store.SetBatchStatus(context.WithoutCancel(r.ctx), r.db, r.batch.ID, biomedb.ImportBatchStatusCanceled); err != nil {
			logrus.Errorf("error setting batch status to canceled for import ID %s: %v", r.batch.ID, err)
		}
		r.withMutex(func() {
			r.status = Cancelled
			r.notify()
		})
		return
	}

	if err := r.store.SetBatchStatus(r.ctx, r.db, r.batch.ID, biomedb.ImportBatchStatusFailed); err != nil {
		logrus.Errorf("error setting batch status to failed for import ID %s: %v", r.batch.ID, err)
	}
	r.withMutex(func() {
		r.cancel()
		r.err = err
		r.status = Failed
	})
}

func (r *ImportRunner) Complete(db db.Querier, userID uuid.UUID) error {
	if err := r.store.SetBatchCompleted(r.ctx, db, r.batch.ID, userID); err != nil {
		r.Fail(err)
		return err
	}
	r.withMutex(func() {
		r.status = Completed
	})
	return nil
}

func (r *ImportRunner) Stop() {
	r.mu.Lock()

	if r.status != Running {
		r.mu.Unlock()
		return
	}

	r.status = Cancelled
	r.cancel()
	r.notify()
	r.mu.Unlock()
}

func (r *ImportRunner) GBIFProgress() progress.ProgressSnapshot {
	return r.gbif.Snapshot()
}

func (r *ImportRunner) notify() {
	r.events.Publish(r.Snapshot())
	r.lastNotify.Store(time.Now())
}

func (r *ImportRunner) Snapshot() BatchSnapshot {
	r.mu.Lock()
	status := r.status
	err := r.err
	r.mu.Unlock()
	return BatchSnapshot{
		ImportID:         r.batch.ID,
		Batch:            r.batch,
		Status:           status,
		ResolutionStatus: r.resolutionStatus,
		GBIF:             r.GBIFProgress(),
		Error:            err,
	}
}

func (r *ImportRunner) TaxonResolver() TaxonResolver {
	return r.taxonResolver
}

func (r *ImportRunner) BibliographyResolver() *BibliographyResolver {
	return r.bibliography
}

func (r *ImportRunner) ResolveTaxon(input models.ResolveInput) error {
	if err := r.taxonResolver.ResolveTaxon(r.ctx, r.db, r.batch.ID, input); err != nil {
		return fmt.Errorf("failed to resolve taxon for import ID %s: %w", r.batch.ID, err)
	}
	r.CheckReadyToMaterialize()
	return nil
}

func (r *ImportRunner) ResolvePublication(input models.ResolveInput) error {
	if err := r.bibliography.ResolvePublication(r.ctx, r.db, r.batch.ID, input); err != nil {
		return fmt.Errorf("failed to resolve publication for import ID %s: %w", r.batch.ID, err)
	}
	r.CheckReadyToMaterialize()
	return nil
}

func (r *ImportRunner) GetMethodsResolution() ([]models.SamplingMethodResolution, error) {
	resolution, err := r.samplings.GetMethodsResolution(r.ctx, r.db, r.batch.ID)
	if err != nil {
		return []models.SamplingMethodResolution{}, err
	}
	return resolution, nil
}

func (r *ImportRunner) ResolveMethod(importID uuid.UUID, input models.SamplingMethodResolutionInput) (models.SamplingMethodResolution, error) {
	return r.samplings.ResolveMethod(r.ctx, r.db, importID, input)
}

func (r *ImportRunner) GetFixativesResolution() ([]models.SamplingFixativeResolution, error) {
	resolution, err := r.samplings.GetFixativesResolution(r.ctx, r.db, r.batch.ID)
	if err != nil {
		return []models.SamplingFixativeResolution{}, err
	}
	return resolution, nil
}

func (r *ImportRunner) ResolveFixative(importID uuid.UUID, input models.SamplingFixativeResolutionInput) (models.SamplingFixativeResolution, error) {
	return r.samplings.ResolveFixative(r.ctx, r.db, importID, input)
}

func (r *ImportRunner) CheckReadyToMaterialize() models.MaterializationReadyCheck {
	check, err := r.store.CheckReadyToMaterialize(r.ctx, r.db, r.batch.ID)
	if err != nil {
		r.Fail(err)
	}
	logrus.Debugf("CheckReadyToMaterialize for import ID %s: %+v", r.batch.ID, check)
	r.withMutex(func() {
		r.resolutionStatus = check
		if check.IsReady() {
			r.status = ReadyToMaterialize
		} else {
			r.status = NeedsResolution
		}
	})
	r.notify()
	return check
}

func (r *ImportRunner) Materialize(userID uuid.UUID) (*models.ImportBatch, error) {
	check := r.CheckReadyToMaterialize()
	if r.Status() != ReadyToMaterialize {
		return nil, fmt.Errorf("cannot materialize batch in status %s : %w", r.status, check.AppError())
	}
	r.withMutex(func() {
		r.status = Materializing
	})
	if err := r.taxonResolver.FillGBIFDependencies(r.ctx, r.db, r.batch.ID, r.gbif); err != nil {
		return nil, fmt.Errorf("fill GBIF dependencies: %w", err)
	}
	if err := r.db.WithTx(r.ctx, func(tx *db.Tx) error {
		logrus.Infof("Materializing import batch: %s", r.batch.Label)

		logrus.Infof("Materializing taxa for import batch: %s", r.batch.Label)
		if err := r.taxonResolver.MaterializeTaxa(r.ctx, tx, r.batch.ID); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				return fmt.Errorf("materialize taxa: %w ; details: %s", err, pgErr.Detail)
			}
			return fmt.Errorf("materialize taxa: %w", err)
		}

		logrus.Infof("Materializing samplings for import batch: %s", r.batch.Label)
		if err := r.samplings.MaterializeSamplings(r.ctx, tx, r.batch.ID); err != nil {
			return fmt.Errorf("materialize samplings: %w", err)
		}

		logrus.Infof("Materializing occurrences for import batch: %s", r.batch.Label)
		if err := r.occurrences.MaterializeOccurrences(r.ctx, tx, r.batch.ID); err != nil {
			return fmt.Errorf("materialize occurrences: %w", err)
		}

		logrus.Infof("Materializing bibliography for import batch: %s", r.batch.Label)
		if err := r.bibliography.MaterializeBibliography(r.ctx, tx, r.batch.ID); err != nil {
			return fmt.Errorf("materialize bibliography: %w", err)
		}

		logrus.Infof("Refreshing occurrence codes for import batch: %s", r.batch.Label)
		if err := r.occurrences.RefreshOccurrenceCodes(r.ctx, tx); err != nil {
			return fmt.Errorf("refresh occurrence codes: %w", err)
		}
		if err := r.Complete(tx, userID); err != nil {
			return fmt.Errorf("complete import: %w", err)
		}
		return nil
	}); err != nil {
		r.Fail(err)
		return nil, err
	}
	return &r.batch, nil
}

func (r *ImportRunner) Delete(ctx context.Context) error {
	if err := r.store.DeleteBatch(ctx, r.db, r.batch.ID); err != nil {
		return fmt.Errorf("delete batch: %w", err)
	}
	r.Stop()
	return nil
}
