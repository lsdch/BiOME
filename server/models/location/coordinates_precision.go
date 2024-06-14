package location

type CoordinatePrecision string

//generate:enum
const (
	M10     CoordinatePrecision = "m10"
	M100    CoordinatePrecision = "m100"
	KM1     CoordinatePrecision = "Km1"
	KM10    CoordinatePrecision = "Km10"
	KM100   CoordinatePrecision = "Km100"
	Unknown CoordinatePrecision = "Unknown"
)
