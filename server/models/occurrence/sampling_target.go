package occurrence

type SamplingTargetKind string

//generate:enum
const (
	Community     SamplingTargetKind = "Community"
	UnknownTarget SamplingTargetKind = "Unknown"
	Taxa          SamplingTargetKind = "Taxa"
)
