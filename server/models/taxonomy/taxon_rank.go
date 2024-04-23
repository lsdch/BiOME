package taxonomy

// TaxonRank represents the taxonomic rank of a taxon.
type TaxonRank string // @name TaxonRank

//generate:enum
const (
	Kingdom    TaxonRank = "Kingdom"
	Phylum     TaxonRank = "Phylum"
	Class      TaxonRank = "Class"
	Family     TaxonRank = "Family"
	Genus      TaxonRank = "Genus"
	Species    TaxonRank = "Species"
	Subspecies TaxonRank = "Subspecies"
)
