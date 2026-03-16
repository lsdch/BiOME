package taxonomy

// TaxonStatus represents the taxonomic status of a taxon.
type TaxonStatus string // @name TaxonStatus

//generate:enum
const (
	Accepted     TaxonStatus = "Accepted"
	Synonym      TaxonStatus = "Synonym"
	Unreferenced TaxonStatus = "Unreferenced"
	Unclassified TaxonStatus = "Unclassified"
)
