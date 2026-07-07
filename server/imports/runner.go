package imports

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/lib/progress"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/stores"
)

const (
	ErrGBIFAlreadyRunning string = "GBIF import is already in progress for this import hash"
)

type RunnerStatus string

const (
	Created   RunnerStatus = "created"
	Running   RunnerStatus = "running"
	Failed    RunnerStatus = "failed"
	Completed RunnerStatus = "completed"
	Cancelled RunnerStatus = "cancelled"
)

type ImportRunner struct {
	db            *db.DB
	workflow      models.ImportWorkflow
	store         *stores.WorkflowStore
	taxonResolver TaxonResolver
	samplings     *services.SamplingService

	parser CSVParser

	status RunnerStatus
	gbif   *progress.ProgressTracker

	mu         sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	err        error
	events     ImportEventSink[ImportEvent]
	lastNotify time.Time
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
	runner.gbif = progress.NewProgressTracker().WithCallback(runner.notify)
	return runner
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

func (r *ImportRunner) StartWorkflowCSV(csvFile io.Reader, separator rune) error {

	if err := r.Start(); err != nil {
		return err
	}

	rows, err := r.parser.ParseCSV(csvFile, separator)
	if err != nil {
		return fmt.Errorf("parse CSV: %w", err)
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
		return nil
	})
	if err != nil {
		r.Fail(err)
		return err
	}

	go r.EnrichGBIF()
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

	r.gbif.Start(int32(len(toFetch)))
	taxa, err := r.taxonResolver.FetchCandidatesFromGBIF(r.ctx, r.db, r.workflow.ImportID, toFetch, r.gbif)
	if err != nil {
		r.Fail(err)
		r.gbif.Fail(err)
		return
	}

	err = r.taxonResolver.InsertGBIFCandidates(r.ctx, r.db, taxa)
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

func (r *ImportRunner) Start() error {
	r.mu.Lock()

	switch r.status {
	case Created:
		r.status = Running
	default:
		r.mu.Unlock()
		return errors.New("workflow cannot be restarted")
	}
	r.mu.Unlock()
	r.notify()
	return nil
}

func (r *ImportRunner) Fail(err error) {
	r.mu.Lock()

	if errors.Is(err, context.Canceled) {
		r.status = Cancelled
		r.mu.Unlock()
		r.notify()
		return
	}

	r.err = err
	r.status = Failed

	r.mu.Unlock()
	r.notify()
}

func (r *ImportRunner) Complete() {
	r.mu.Lock()
	r.status = Completed
	r.mu.Unlock()
	r.notify()
}

func (r *ImportRunner) Stop() {
	r.mu.Lock()

	if r.status != Running {
		r.mu.Unlock()
		return
	}

	r.status = Cancelled
	r.mu.Unlock()

	r.cancel()
	r.notify()
}

func (r *ImportRunner) GBIFProgress() progress.ProgressSnapshot {
	return r.gbif.Snapshot()
}

func (r *ImportRunner) notify() {
	if time.Since(r.lastNotify) < 1000*time.Millisecond {
		return
	}
	r.lastNotify = time.Now()
	r.events.Publish(r.Snapshot())
}

func (r *ImportRunner) Snapshot() ImportEvent {
	r.mu.Lock()
	status := r.status
	err := r.err
	r.mu.Unlock()
	return ImportEvent{
		Workflow: r.workflow,
		Status:   status,
		GBIF:     r.GBIFProgress(),
		Error:    err,
	}
}
