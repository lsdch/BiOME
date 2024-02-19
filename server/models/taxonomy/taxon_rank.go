package taxonomy

// TaxonRank represents the taxonomic rank of a taxon.
type TaxonRank string // @name TaxonRank

const (
	Kingdom    TaxonRank = "Kingdom"
	Phylum     TaxonRank = "Phylum"
	Class      TaxonRank = "Class"
	Family     TaxonRank = "Family"
	Genus      TaxonRank = "Genus"
	Species    TaxonRank = "Species"
	Subspecies TaxonRank = "Subspecies"
)

func (m *TaxonRank) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(*m), nil
}

func (m *TaxonRank) UnmarshalEdgeDBStr(data []byte) error {
	*m = TaxonRank(string(data))
	return nil
}
