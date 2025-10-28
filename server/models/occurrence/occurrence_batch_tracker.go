package occurrence

import (
	"github.com/schollz/progressbar/v3"
)

type OccurrenceBatchTracker interface {
	Start(total int) OccurrenceBatchTracker
	StartUnknownTotal(description string) OccurrenceBatchTracker
	SetDescription(description string) OccurrenceBatchTracker
	Progress(quantity int) OccurrenceBatchTracker
	SetDetail(detail string) OccurrenceBatchTracker
	SetDetailsPrefix(prefix string) OccurrenceBatchTracker
	Finish() OccurrenceBatchTracker
}

type NoOpBatchTracker struct{}

func (o *NoOpBatchTracker) Start(total int) OccurrenceBatchTracker {
	return o
}

func (o *NoOpBatchTracker) StartUnknownTotal(description string) OccurrenceBatchTracker {
	return o
}

func (o *NoOpBatchTracker) SetDescription(description string) OccurrenceBatchTracker {
	return o
}

func (o *NoOpBatchTracker) Progress(quantity int) OccurrenceBatchTracker {
	return o
}

func (o *NoOpBatchTracker) SetDetail(detail string) OccurrenceBatchTracker {
	return o
}

func (o *NoOpBatchTracker) SetDetailsPrefix(prefix string) OccurrenceBatchTracker {
	return o
}

func (o *NoOpBatchTracker) Finish() OccurrenceBatchTracker {
	return o
}

var _ OccurrenceBatchTracker = (*NoOpBatchTracker)(nil)

type OccurrenceBatchProgressBar struct {
	pb            *progressbar.ProgressBar
	detailsPrefix string
}

func (o *OccurrenceBatchProgressBar) StartUnknownTotal(description string) OccurrenceBatchTracker {
	o.pb = progressbar.Default(-1, description)
	_ = o.pb.RenderBlank()
	return o
}

func (o *OccurrenceBatchProgressBar) SetDescription(description string) OccurrenceBatchTracker {
	o.pb.Describe(description)
	return o
}

func (o *OccurrenceBatchProgressBar) Start(total int) OccurrenceBatchTracker {
	o.pb = progressbar.NewOptions(
		total,
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetMaxDetailRow(5),
	)
	return o
}

func (o *OccurrenceBatchProgressBar) SetDetail(detail string) OccurrenceBatchTracker {
	_ = o.pb.AddDetail(o.detailsPrefix + detail)
	return o
}
func (o *OccurrenceBatchProgressBar) SetDetailsPrefix(prefix string) OccurrenceBatchTracker {
	o.detailsPrefix = prefix
	return o
}

func (o *OccurrenceBatchProgressBar) Progress(quantity int) OccurrenceBatchTracker {
	_ = o.pb.Add(quantity)
	return o
}

func (o *OccurrenceBatchProgressBar) Finish() OccurrenceBatchTracker {
	_ = o.pb.Clear()
	_ = o.pb.Close()
	return o
}

var _ OccurrenceBatchTracker = (*OccurrenceBatchProgressBar)(nil)
