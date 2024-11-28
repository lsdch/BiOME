package occurrence

type DatePrecision string

//generate:enum
const (
	Day     DatePrecision = "Day"
	Month   DatePrecision = "Month"
	Year    DatePrecision = "Year"
	Unknown DatePrecision = "Unknown"
)
