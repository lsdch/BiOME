package occurrence

type DatePrecision string

//generate:enum
const (
	Year    DatePrecision = "Year"
	Month   DatePrecision = "Month"
	Day     DatePrecision = "Day"
	Unknown DatePrecision = "Unknown"
)
