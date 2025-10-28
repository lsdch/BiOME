package specimen

type Quantity string

// generate:enum
const (
	One     Quantity = "One"
	Several Quantity = "Several"
	Dozen   Quantity = "Dozen"
	Tens    Quantity = "Tens"
	Hundred Quantity = "Hundred"
)
