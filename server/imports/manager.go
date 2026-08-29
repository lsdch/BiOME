package imports

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/services/storage"
	"github.com/lsdch/biome/stores"
)

type ImportManager struct {
	db *db.DB
	mu sync.RWMutex

	batchStore  *stores.BatchesStore
	fileStorage storage.RawFileStorage

	taxonResolver TaxonResolver
	bibliography  *BibliographyResolver

	samplings   *services.SamplingService
	occurrences *services.OccurrencesService

	runners map[uuid.UUID]*ImportRunner
	broker  *EventBroker[BatchSnapshot]
}

func NewImportManager(db *db.DB,
	batchStore *stores.BatchesStore,
	taxonResolver TaxonResolver,
	bibliography *BibliographyResolver,
	samplings *services.SamplingService,
	occurrences *services.OccurrencesService,
	storage storage.RawFileStorage,
) *ImportManager {
	return &ImportManager{
		mu:            sync.RWMutex{},
		db:            db,
		batchStore:    batchStore,
		taxonResolver: taxonResolver,
		bibliography:  bibliography,
		samplings:     samplings,
		occurrences:   occurrences,
		fileStorage:   storage,
		runners:       make(map[uuid.UUID]*ImportRunner),
		broker:        NewEventBroker[BatchSnapshot](),
	}
}

func (m ImportManager) FileStorage() storage.RawFileStorage {
	return m.fileStorage
}

func (m *ImportManager) addRunner(runner *ImportRunner) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.runners[runner.Batch().ID] = runner
}

func (m *ImportManager) NewBatch(ctx context.Context, userID uuid.UUID, w models.ImportBatchWithFileInput) (*ImportRunner, error) {
	batch, err := m.batchStore.CreateBatch(ctx, m.db, userID, w)
	if err != nil {
		return nil, err
	}
	runner := NewImportRunner(context.Background(),
		m.db, m.broker,
		batch, m.batchStore,
		m.samplings, m.taxonResolver, m.bibliography,
		m.occurrences,
		m.fileStorage)
	m.addRunner(runner)
	return runner, nil
}

func (m *ImportManager) Restore(ctx context.Context) error {
	batchs, err := m.batchStore.ListBatches(ctx, m.db)
	if err != nil {
		return err
	}
	for _, batch := range batchs {
		if batch.Status == biomedb.ImportBatchStatusCompleted {
			continue
		}
		runner := NewImportRunner(ctx,
			m.db, m.broker,
			batch, m.batchStore,
			m.samplings, m.taxonResolver, m.bibliography,
			m.occurrences,
			m.fileStorage)
		runner.Run()
		m.addRunner(runner)
	}
	return nil
}

func (m *ImportManager) GetRunner(importID uuid.UUID) (*ImportRunner, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	runner, ok := m.runners[importID]
	return runner, ok
}

func (m *ImportManager) ListRunners() []*ImportRunner {
	m.mu.RLock()
	defer m.mu.RUnlock()
	runners := make([]*ImportRunner, 0, len(m.runners))
	for _, runner := range m.runners {
		runners = append(runners, runner)
	}
	return runners
}

// RemoveRunner stops and removes the ImportRunner associated with the given importID from the manager.
//
// It does not delete the import batch from the database; it only stops the runner and removes it from the manager's internal map.
func (m *ImportManager) RemoveRunner(importID uuid.UUID) {
	m.mu.Lock()
	runner, ok := m.runners[importID]
	delete(m.runners, importID)
	m.mu.Unlock()

	if ok {
		runner.Stop()
	}
}

func (m *ImportManager) Subscribe() (events <-chan BatchSnapshot, unsubscribe func()) {
	return m.broker.Subscribe()
}

func (m *ImportManager) Snapshots() []BatchSnapshot {
	snapshots := make([]BatchSnapshot, len(m.runners))
	for _, runner := range m.runners {
		snapshots = append(snapshots, runner.Snapshot())
	}
	return snapshots
}

func (m *ImportManager) SnapshotsForUser(userID uuid.UUID) []BatchSnapshot {
	snapshots := make([]BatchSnapshot, 0, len(m.runners))
	for _, runner := range m.runners {
		snapshot := runner.Snapshot()
		if snapshot.Batch.CreatedBy == userID {
			snapshots = append(snapshots, snapshot)
		}
	}
	return snapshots
}
