package gbif

import (
	"darco/proto/models/taxonomy"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// ProgressTracker tracks and reports the progress of the import process.
type ProgressTracker struct {
	process ImportProcess
	monitor func(p *ImportProcess)
}

func (p *ProgressTracker) Report() {
	p.monitor(&p.process)
}

func (p *ProgressTracker) Progress(n int) {
	p.process.Imported += n
	logrus.Debugf("Progress")
	p.Report()
}

func (p *ProgressTracker) Errorf(format string, a ...any) error {
	p.process.Error = fmt.Errorf(format, a...)
	logrus.Errorf(format, a...)
	p.Terminate()
	return p.process.Error
}

func (p *ProgressTracker) Terminate() {
	p.process.Done = true
	p.Report()
}

func NewProgressTracker(taxon *TaxonGBIF, f func(p *ImportProcess)) *ProgressTracker {
	process := ImportProcess{
		Name:     taxon.Name,
		GBIF_ID:  taxon.Key,
		Expected: taxon.NumDescendants + 1,
		Imported: 0,
		Rank:     taxonomy.TaxonRank(taxon.Rank),
		Started:  time.Now(),
		Done:     false,
		Error:    nil,
	}
	tracker := ProgressTracker{
		process: process,
		monitor: f,
	}
	tracker.Report()
	return &tracker
}
