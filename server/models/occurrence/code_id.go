package occurrence

import "time"

type CodeHistory struct {
	Code string    `edgedb:"code" json:"code"`
	Time time.Time `edgedb:"time" json:"time"`
}

type CodeIdentifier struct {
	Code        string        `edgedb:"code" json:"code"`
	CodeHistory []CodeHistory `edgedb:"code_history" json:"code_history,omitempty"`
}
