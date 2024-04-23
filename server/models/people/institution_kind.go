package people

type InstitutionKind string // @name InstitutionKind

//generate:enum
const (
	Lab                InstitutionKind = "Lab"
	FoundingAgency     InstitutionKind = "FundingAgency"
	SequencingPlatform InstitutionKind = "SequencingPlatform"
	Other              InstitutionKind = "Other"
)
