package occurrence

import (
	"time"

	"github.com/geldata/gel-go/geltypes"
)

type CodeHistory struct {
	Code string    `gel:"code" json:"code"`
	Time time.Time `gel:"time" json:"time"`
}

type CodeIdentifier struct {
	Code        string        `gel:"code" json:"code"`
	CodeHistory []CodeHistory `gel:"code_history" json:"code_history,omitempty"`
}

type CreatedCode struct {
	ID geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	// Code is a unique identifier for the occurrence within the system.
	Code string `gel:"code" json:"code"`
}
