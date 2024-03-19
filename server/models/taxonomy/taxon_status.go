package taxonomy

// TaxonStatus represents the taxonomic status of a taxon.
type TaxonStatus string // @name TaxonStatus

const (
	Accepted     TaxonStatus = "Accepted"
	Synonym      TaxonStatus = "Synonym"
	Unclassified TaxonStatus = "Unclassified"
)

func (m TaxonStatus) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *TaxonStatus) UnmarshalEdgeDBStr(data []byte) error {
	*m = TaxonStatus(string(data))
	return nil
}
