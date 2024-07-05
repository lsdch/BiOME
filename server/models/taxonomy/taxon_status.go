package taxonomy

// TaxonStatus represents the taxonomic status of a taxon.
type TaxonStatus string // @name TaxonStatus

//generate:enum
const (
	Accepted     TaxonStatus = "Accepted"
	Unreferenced TaxonStatus = "Unreferenced"
	Unclassified TaxonStatus = "Unclassified"
)
