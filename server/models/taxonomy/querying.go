package taxonomy

import (
	"context"

	"github.com/geldata/gel-go/geltypes"
)

type TaxaFilters struct {
	Taxa       []string               `query:"taxa" json:"taxa,omitempty"`
	WholeClade bool                   `query:"whole_clade" json:"whole_clade,omitempty" default:"false"`
	TaxaByRank map[TaxonRank][]string `json:"-"`
}

func (i *TaxaFilters) FetchTaxa(db geltypes.Executor) error {
	if len(i.Taxa) > 0 && i.WholeClade {
		i.TaxaByRank = make(map[TaxonRank][]string)
		var allTaxa []TaxonInner
		err := db.Query(context.Background(),
			`#edgeql
			with module taxonomy
			select distinct taxonomy::Taxon { name, rank }
			filter .name in array_unpack(<array<str>>$0)
			`,
			&allTaxa, i.Taxa)
		if err != nil {
			return err
		}
		for _, taxon := range allTaxa {
			i.TaxaByRank[taxon.Rank] = append(i.TaxaByRank[taxon.Rank], taxon.Name)
		}
	}
	return nil
}
