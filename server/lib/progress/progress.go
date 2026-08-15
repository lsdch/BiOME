package progress

import (
	"sync"
	"time"

	"github.com/lsdch/biome/models"
)

type ProgressReporter interface {
	Start(total int32)
	Increment(int32)
	Complete()
	Fail(error)
	AddToTotal(int32)
}

type NoopProgressReporter struct{}

func (n *NoopProgressReporter) Start(int32)      {}
func (n *NoopProgressReporter) Increment(int32)  {}
func (n *NoopProgressReporter) Complete()        {}
func (n *NoopProgressReporter) Fail(error)       {}
func (n *NoopProgressReporter) AddToTotal(int32) {}

type ProgressStatus string

const (
	Pending   ProgressStatus = "pending"
	Running   ProgressStatus = "running"
	Completed ProgressStatus = "completed"
	Failed    ProgressStatus = "failed"
)

type ProgressTracker struct {
	mu sync.RWMutex

	status      ProgressStatus
	total       int32
	completed   int32
	startedAt   models.Optional[time.Time]
	completedAt models.Optional[time.Time]
	errorMsg    string

	onUpdate func(status ProgressStatus)
}

func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		status: Pending,
	}
}

func (g *ProgressTracker) WithCallback(onUpdate func(status ProgressStatus)) *ProgressTracker {
	g.onUpdate = onUpdate
	return g
}

func (g *ProgressTracker) notify() {
	if g.onUpdate != nil {
		g.onUpdate(g.status)
	}
}

func (g *ProgressTracker) Increment(n int32) {
	g.mu.Lock()
	g.completed += n
	g.mu.Unlock()
	g.notify()
}

func (g *ProgressTracker) Start(total int32) {
	g.mu.Lock()

	g.status = Running
	g.startedAt = models.NewOptional(time.Now())
	g.completedAt.Clear()
	g.total = total
	g.errorMsg = ""

	g.mu.Unlock()
	g.notify()
}

func (g *ProgressTracker) AddToTotal(n int32) {
	g.mu.Lock()
	g.total += n
	g.mu.Unlock()
	g.notify()
}

func (g *ProgressTracker) Complete() {
	g.mu.Lock()

	g.status = Completed
	g.completedAt = models.NewOptional(time.Now())
	g.completed = g.total

	g.mu.Unlock()
	g.notify()
}

func (g *ProgressTracker) Fail(err error) {
	g.mu.Lock()

	g.status = Failed
	g.completedAt.Clear()
	if err != nil {
		g.errorMsg = err.Error()
	}

	g.mu.Unlock()
	g.notify()
}

type ProgressSnapshot struct {
	Status      ProgressStatus
	Total       int32
	Completed   int32
	StartedAt   models.Optional[time.Time]
	CompletedAt models.Optional[time.Time]
	Error       string
}

func (p *ProgressTracker) Snapshot() ProgressSnapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return ProgressSnapshot{
		Status:      p.status,
		Total:       p.total,
		Completed:   p.completed,
		StartedAt:   p.startedAt,
		CompletedAt: p.completedAt,
		Error:       p.errorMsg,
	}
}
