package occurrence

import (
	"context"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/dataset"
)

type SequenceDatasetListItem struct {
	dataset.DatasetInner `gel:"$inline" json:",inline"`
	Sites                int64 `gel:"sites_count" json:"sites"`
	Sequences            int64 `gel:"sequences_count" json:"sequences"`
}

type SequenceDataset struct {
	dataset.Dataset `gel:"$inline" json:",inline"`
	Sites           []SiteItem `gel:"sites" json:"sites"`
	Sequences       []Sequence `gel:"sequences" json:"sequences"`
}

func ListSequenceDatasets(db geltypes.Executor) (datasets []SequenceDatasetListItem, err error) {
	err = db.Query(context.Background(),
		`#edgeql
			select datasets::SeqDataset {
				*,
				sites_count := count(.sites),
				sequences_count := count(.sequences),
			}
		`,
		&datasets,
	)
	return
}

func GetSequenceDataset(db geltypes.Executor, slug string) (dataset SequenceDataset, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select datasets::SequenceDataset {
				**,
				sites: { *, country: { * } },
				sequences: {
					**,
					gene: { * },
					sampling { *, site: { *, country: { * } } },
					identification: { ** },
					external := [is seq::ExternalSequence]{
						origin,
						referenced_in: { ** },
						published_in: { ** },
						specimen_identifier,
						original_taxon,
						source_sample : {
							[is occurrence::BioMaterial].*,
							identification: { ** }
						}
					}
				}
			} filter .slug = <str>$0
		`,
		&dataset, slug,
	)
	return
}
