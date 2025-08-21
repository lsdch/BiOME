package taxonomy

// TaxonRank represents the taxonomic rank of a taxon.
type TaxonRank string

//generate:enum
const (
	Kingdom    TaxonRank = "Kingdom"
	Phylum     TaxonRank = "Phylum"
	Class      TaxonRank = "Class"
	Order      TaxonRank = "Order"
	Family     TaxonRank = "Family"
	Genus      TaxonRank = "Genus"
	Subgenus   TaxonRank = "Subgenus"
	Species    TaxonRank = "Species"
	Subspecies TaxonRank = "Subspecies"
)
