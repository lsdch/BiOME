package sequences

type ExtSeqOrigin string

//generate:enum
const (
	Lab     ExtSeqOrigin = "Lab"
	DB      ExtSeqOrigin = "DB"
	PersCom ExtSeqOrigin = "PersCom"
)
