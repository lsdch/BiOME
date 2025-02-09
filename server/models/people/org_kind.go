package people

type OrgKind string

//generate:enum
const (
	Lab                OrgKind = "Lab"
	FoundingAgency     OrgKind = "FundingAgency"
	SequencingPlatform OrgKind = "SequencingPlatform"
	Other              OrgKind = "Other"
)
