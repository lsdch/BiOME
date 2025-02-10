package occurrence

import (
	"context"

	"github.com/edgedb/edgedb-go"
	"github.com/lsdch/biome/models/dataset"
)

type SequenceDataset struct {
	dataset.AbstractDataset `edgedb:"$inline" json:",inline"`
	Sites                   []SiteItem `edgedb:"sites" json:"sites"`
	Sequences               []Sequence `edgedb:"sequences" json:"sequences"`
}

func ListSequenceDatasets(db edgedb.Executor) (datasets []SequenceDataset, err error) {
	err = db.Query(context.Background(),
		`#edgeql
			select datasets::SeqDataset {
				**,
				sites: { *, country: { * } },
				sequences: {
					**,
					gene: { * },
					required event := .sampling.event { *, site: {name, code} },
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
			}
		`,
		&datasets,
	)
	return
}

func GetSequenceDataset(db edgedb.Executor, slug string) (dataset SequenceDataset, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select datasets::SequenceDataset {
				**,
				sites: { *, country: { * } },
				sequences: {
					**,
					gene: { * },
					required event := .sampling.event { *, site: {name, code} },
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
