package imports

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/stores"
)

type ImportManager struct {
	db            *db.DB
	mu            sync.RWMutex
	workflowStore *stores.WorkflowStore
	taxonResolver TaxonResolver
	samplings     *services.SamplingService
	runners       map[uuid.UUID]*ImportRunner
	broker        *EventBroker[ImportEvent]
}

func NewImportManager(db *db.DB, workflowStore *stores.WorkflowStore, taxonResolver TaxonResolver, samplings *services.SamplingService) *ImportManager {
	return &ImportManager{
		mu:            sync.RWMutex{},
		db:            db,
		workflowStore: workflowStore,
		taxonResolver: taxonResolver,
		samplings:     samplings,
		runners:       make(map[uuid.UUID]*ImportRunner),
		broker:        NewEventBroker[ImportEvent](),
	}
}

func (m *ImportManager) addRunner(runner *ImportRunner) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.runners[runner.Workflow().ImportID] = runner
}

func (m *ImportManager) NewWorkflow(ctx context.Context, label string) (*ImportRunner, error) {
	workflow, err := m.workflowStore.CreateWorkflow(ctx, m.db, label)
	if err != nil {
		return nil, err
	}
	runner := NewImportRunner(ctx, m.db, m.broker, workflow, m.workflowStore, m.samplings, m.taxonResolver)
	m.addRunner(runner)
	return runner, nil
}

func (m *ImportManager) Restore(ctx context.Context) error {
	workflows, err := m.workflowStore.ListWorkflows(ctx, m.db)
	if err != nil {
		return err
	}
	for _, workflow := range workflows {
		runner := NewImportRunner(ctx, m.db, m.broker, workflow, m.workflowStore, m.samplings, m.taxonResolver)
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

func (m *ImportManager) RemoveRunner(importID uuid.UUID) {
	m.mu.Lock()
	runner, ok := m.runners[importID]
	delete(m.runners, importID)
	m.mu.Unlock()

	if ok {
		runner.Stop()
	}
}

func (m *ImportManager) Subscribe() (events <-chan ImportEvent, unsubscribe func()) {
	return m.broker.Subscribe()
}

func (m *ImportManager) Snapshots() []ImportEvent {
	snapshots := make([]ImportEvent, len(m.runners))
	for _, runner := range m.runners {
		snapshots = append(snapshots, runner.Snapshot())
	}
	return snapshots
}
