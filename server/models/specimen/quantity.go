package specimen

type Quantity string

// generate:enum
const (
	Unknown Quantity = "Unknown"
	One     Quantity = "One"
	Several Quantity = "Several"
	Dozen   Quantity = "Dozen"
	Tens    Quantity = "Tens"
	Hundred Quantity = "Hundred"
)
