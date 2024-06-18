package location

type CoordinatePrecision string

//generate:enum
const (
	M100    CoordinatePrecision = "<100m"
	KM1     CoordinatePrecision = "<1km"
	KM10    CoordinatePrecision = "<10km"
	KM100   CoordinatePrecision = "10-100km"
	Unknown CoordinatePrecision = "Unknown"
)
