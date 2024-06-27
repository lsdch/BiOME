package location

type CoordinatesPrecision string

//generate:enum
const (
	M100    CoordinatesPrecision = "<100m"
	KM1     CoordinatesPrecision = "<1km"
	KM10    CoordinatesPrecision = "<10km"
	KM100   CoordinatesPrecision = "10-100km"
	Unknown CoordinatesPrecision = "Unknown"
)
