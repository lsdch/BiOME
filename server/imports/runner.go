package imports

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/lib/progress"
	"github.com/lsdch/biome/models"
	csvmodels "github.com/lsdch/biome/models/csv"
	"github.com/lsdch/biome/services"
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
	workflow      models.ImportWorkflow
	store         *stores.WorkflowStore
	taxonResolver TaxonResolver
	samplings     *services.SamplingService
	occurrences   *services.OccurrencesService

	parser CSVParser

	status           RunnerStatus
	resolutionStatus models.MaterializationReadyCheck
	gbif             *progress.ProgressTracker

	mu         sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	err        error
	events     ImportEventSink[ImportEvent]
	lastNotify atomic.Value //time.Time
}

func NewImportRunner(ctx context.Context, db *db.DB,
	events ImportEventSink[ImportEvent],
	workflow models.ImportWorkflow,
	store *stores.WorkflowStore,
	samplings *services.SamplingService,
	taxonResolver TaxonResolver,
) *ImportRunner {
	ctx, cancel := context.WithCancel(ctx)

	runner := &ImportRunner{
		db:            db,
		workflow:      workflow,
		store:         store,
		taxonResolver: taxonResolver,
		samplings:     samplings,
		status:        Created,
		parser:        NewCSVParser(),
		ctx:           ctx,
		cancel:        cancel,
		events:        events,
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
	return r.workflow.ImportID
}

func (r *ImportRunner) Workflow() models.ImportWorkflow {
	return r.workflow
}

func (r *ImportRunner) Status() RunnerStatus {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.status
}

type RowValidationErrors struct {
	RowNumber int32
	Errors    validator.ValidationErrors
}

func (e *RowValidationErrors) Error() string {
	return fmt.Sprintf("row %d: %v", e.RowNumber, e.Errors)
}

func (e *RowValidationErrors) Unwrap() error {
	return e.Errors
}

type BatchValidationErrors struct {
	Errors []*RowValidationErrors
}

func (e *BatchValidationErrors) Error() string {
	var sb strings.Builder
	sb.WriteString("validation errors:\n")
	for _, rowErr := range e.Errors {
		sb.WriteString(rowErr.Error())
	}
	return sb.String()
}

func (e *BatchValidationErrors) Unwrap() []error {
	errs := make([]error, len(e.Errors))

	for i, rowErr := range e.Errors {
		errs[i] = rowErr
	}

	return errs
}

func (r *ImportRunner) ValidateRows(rows []csvmodels.OccurrenceImportRow) error {
	validationErrors := BatchValidationErrors{Errors: make([]*RowValidationErrors, 0, len(rows))}
	for rowNum, row := range rows {
		if err := row.Validate(); err != nil {
			var validationErrs validator.ValidationErrors
			if errors.As(err, &validationErrs) {
				validationErrors.Errors = append(validationErrors.Errors, &RowValidationErrors{
					RowNumber: int32(rowNum + 2),
					Errors:    validationErrs,
				})
			}
		}
	}
	if len(validationErrors.Errors) > 0 {
		return &validationErrors
	}
	return nil
}

func (r *ImportRunner) StartWorkflowCSV(csvFile io.Reader, separator rune) error {

	if err := r.StartStaging(); err != nil {
		return err
	}

	rows, err := r.parser.ParseCSV(csvFile, separator)
	if err != nil {
		return err
	}

	if err := r.ValidateRows(rows); err != nil {
		r.Fail(err)
		return err
	}

	err = r.db.WithTx(r.ctx, func(tx *db.Tx) error {
		if err := r.store.InsertStaging(r.ctx, tx, r.workflow.ImportID, rows); err != nil {
			return fmt.Errorf("insert staging occurrences: %w", err)
		}
		if _, err = r.taxonResolver.InitResolution(r.ctx, tx, r.workflow.ImportID); err != nil {
			return fmt.Errorf("init taxon resolution: %w", err)
		}
		if _, err := r.samplings.InitMethodResolution(r.ctx, tx, r.workflow.ImportID); err != nil {
			return fmt.Errorf("init method resolution: %w", err)
		}
		if _, err := r.samplings.InitFixativeResolution(r.ctx, tx, r.workflow.ImportID); err != nil {
			return fmt.Errorf("init fixative resolution: %w", err)
		}
		return nil
	})
	if err != nil {
		r.Fail(err)
		return err
	}

	r.MarkStaged()

	go r.Run()

	return nil
}

func (r *ImportRunner) Runnable() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.status == Staged || r.status == Failed || r.status == Cancelled
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
	}
	wg.Wait()
	logrus.Infof("Running automatic taxon resolution for unambiguous candidates for import ID %s", r.workflow.ImportID)
	if err := r.taxonResolver.AutoResolveUnambiguousCandidates(r.ctx, r.db, r.workflow.ImportID); err != nil {
		logrus.Errorf("autoresolve failed for import ID %s: %v", r.workflow.ImportID, err)
		r.Fail(fmt.Errorf("autoresolve failed: %w", err))
		return err
	}

	r.CheckReadyToMaterialize()
	return nil
}

func (r *ImportRunner) EnrichGBIF() {

	if r.gbif.Snapshot().Status == progress.Running {
		return
	}

	toFetch, err := r.taxonResolver.ListTaxaToFetch(r.ctx, r.db, r.workflow.ImportID)
	if err != nil {
		r.Fail(err)
		r.gbif.Fail(err)
		return
	} else if len(toFetch) == 0 {
		r.gbif.Complete()
		return
	}

	logrus.Infof("Fetching %d taxa from GBIF for import ID %s", len(toFetch), r.workflow.ImportID)

	r.gbif.Start(int32(len(toFetch)))
	taxa, err := r.taxonResolver.FetchCandidatesFromGBIF(r.ctx, r.db, r.workflow.ImportID, toFetch, r.gbif)
	if err != nil {
		r.Fail(err)
		r.gbif.Fail(err)
		return
	}

	err = r.taxonResolver.InsertGBIFCandidates(r.ctx, r.db, r.workflow.ImportID, taxa)
	if err != nil {
		r.Fail(err)
		r.gbif.Fail(err)
		return
	}

	err = r.taxonResolver.MarkTaxaGBIFImportCompleted(r.ctx, r.db, r.workflow.ImportID, slices.Collect(maps.Keys(taxa)))
	if err != nil {
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
		return fmt.Errorf("cannot start workflow in status %s", r.status)
	}
	return nil
}

func (r *ImportRunner) MarkStaged() {
	r.withMutex(func() {
		r.status = Staged
	})
}

func (r *ImportRunner) Fail(err error) {

	if errors.Is(err, context.Canceled) {
		r.withMutex(func() {
			r.status = Cancelled
			r.notify()
		})
		return
	}
	r.withMutex(func() {
		r.err = err
		r.status = Failed
	})
}

func (r *ImportRunner) Complete() {
	r.withMutex(func() {
		r.status = Completed
	})
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

func (r *ImportRunner) Snapshot() ImportEvent {
	r.mu.Lock()
	status := r.status
	err := r.err
	r.mu.Unlock()
	return ImportEvent{
		Workflow:         r.workflow,
		Status:           status,
		ResolutionStatus: r.resolutionStatus,
		GBIF:             r.GBIFProgress(),
		Error:            err,
	}
}

func (r *ImportRunner) TaxonResolver() TaxonResolver {
	return r.taxonResolver
}

func (r *ImportRunner) GetMethodsResolution() ([]models.SamplingMethodResolution, error) {
	resolution, err := r.samplings.GetMethodsResolution(r.ctx, r.db, r.workflow.ImportID)
	if err != nil {
		return []models.SamplingMethodResolution{}, err
	}
	return resolution, nil
}

func (r *ImportRunner) ResolveMethod(importID uuid.UUID, input models.SamplingMethodResolutionInput) (models.SamplingMethodResolution, error) {
	return r.samplings.ResolveMethod(r.ctx, r.db, importID, input)
}

func (r *ImportRunner) GetFixativesResolution() ([]models.SamplingFixativeResolution, error) {
	resolution, err := r.samplings.GetFixativesResolution(r.ctx, r.db, r.workflow.ImportID)
	if err != nil {
		return []models.SamplingFixativeResolution{}, err
	}
	return resolution, nil
}

func (r *ImportRunner) ResolveFixative(importID uuid.UUID, input models.SamplingFixativeResolutionInput) (models.SamplingFixativeResolution, error) {
	return r.samplings.ResolveFixative(r.ctx, r.db, importID, input)
}

func (r *ImportRunner) CheckReadyToMaterialize() {
	check, err := r.store.CheckReadyToMaterialize(r.ctx, r.db, r.workflow.ImportID)
	if err != nil {
		r.Fail(err)
	}
	r.withMutex(func() {
		r.resolutionStatus = check
		if check.IsReady() {
			r.status = ReadyToMaterialize
		} else {
			r.status = NeedsResolution
		}
	})
	r.notify()
}

func (r *ImportRunner) Materialize() (*models.ImportBatch, error) {
	r.CheckReadyToMaterialize()
	if r.Status() != ReadyToMaterialize {
		return nil, fmt.Errorf("cannot materialize workflow in status %s", r.status)
	}
	r.withMutex(func() {
		r.status = Materializing
	})
	if err := r.taxonResolver.FillGBIFDependencies(r.ctx, r.db, r.workflow.ImportID); err != nil {
		return nil, fmt.Errorf("fill GBIF dependencies: %w", err)
	}
	var createdBatch models.ImportBatch
	if err := r.db.WithTx(r.ctx, func(tx *db.Tx) error {
		batch, err := r.store.MaterializeStaging(r.ctx, tx, r.workflow.ImportID)
		if err != nil {
			return fmt.Errorf("materialize batch: %w", err)
		}
		if err := r.samplings.MaterializeSamplings(r.ctx, tx, batch.ID); err != nil {
			return fmt.Errorf("materialize samplings: %w", err)
		}

		if err := r.taxonResolver.MaterializeTaxa(r.ctx, tx, r.workflow.ImportID); err != nil {
			return fmt.Errorf("materialize taxa: %w", err)
		}
		if err := r.occurrences.MaterializeOccurrences(r.ctx, tx, batch.ID); err != nil {
			return fmt.Errorf("materialize occurrences: %w", err)
		}
		if err := r.occurrences.RefreshOccurrenceCodes(r.ctx, tx); err != nil {
			return fmt.Errorf("refresh occurrence codes: %w", err)
		}
		createdBatch = batch
		return nil
	}); err != nil {
		return nil, err
	}
	r.Complete()
	return &createdBatch, nil
}

func (r *ImportRunner) Delete() error {
	if err := r.store.DeleteWorkflow(r.ctx, r.db, r.workflow.ImportID); err != nil {
		return fmt.Errorf("delete workflow: %w", err)
	}
	r.Stop()
	return nil
}
