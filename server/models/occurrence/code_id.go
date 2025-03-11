package occurrence

import "time"

type CodeHistory struct {
	Code string    `gel:"code" json:"code"`
	Time time.Time `gel:"time" json:"time"`
}

type CodeIdentifier struct {
	Code        string        `gel:"code" json:"code"`
	CodeHistory []CodeHistory `gel:"code_history" json:"code_history,omitempty"`
}
