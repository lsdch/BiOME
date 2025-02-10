package dataset

type DatasetCategory string

//generate:enum
const (
	Site       DatasetCategory = "Site"
	Occurrence DatasetCategory = "Occurrence"
	Seq        DatasetCategory = "Seq"
)
