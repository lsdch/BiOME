package occurrence

import (
	"context"
	"encoding/json"

	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/dataset"

	"github.com/edgedb/edgedb-go"
)

type OccurrenceDataset struct {
	dataset.AbstractDataset `edgedb:"$inline" json:",inline"`
	Sites                   []SiteItem               `edgedb:"sites" json:"sites"`
	Occurrences             []OccurrenceWithCategory `edgedb:"occurrences" json:"occurrences"`
	IsCongruent             bool                     `edgedb:"is_congruent" json:"is_congruent"`
}

func ListOccurrenceDatasets(db edgedb.Executor) ([]OccurrenceDataset, error) {
	datasets := []OccurrenceDataset{}
	err := db.Query(context.Background(),
		`#edgeql
			select datasets::OccurrenceDataset {
				**,
				sites: { *, country: { * } },
				occurrences: {
					sampling: { * },
					identification: { ** },
					comments
				},
			}
		`,
		&datasets,
	)
	return datasets, err
}

func GetOccurrenceDataset(db edgedb.Executor, slug string) (dataset OccurrenceDataset, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select datasets::OccurrenceDataset {
				**,
				sites: { *, country: { * } },
				occurrences: {
					sampling: { * },
					identification: { ** },
					comments
				},
			} filter .slug = <str>$0
		`,
		&dataset, slug,
	)
	return
}

type OccurrenceDatasetInput struct {
	dataset.DatasetInput `edgedb:"$inline" json:",inline"`
	Occurrences          OccurrenceBatchInput `json:"occurrences"`
}

func (i OccurrenceDatasetInput) SaveTx(tx *edgedb.Tx) (created OccurrenceDataset, err error) {
	occurrences, err := i.Occurrences.Save(tx)
	if err != nil {
		return created, models.WrapErrorPath(err, "occurrences")
	}

	i.DatasetInput.GenerateSlug()
	data, _ := json.Marshal(i.DatasetInput)
	occurrencesData, _ := json.Marshal(occurrences)

	err = tx.QuerySingle(context.Background(),
		`#edgeql
      with
        data := <json>$0,
        occurrences := <json>$1,
      select (insert datasets::OccurrenceDataset {
        label := <str>data['label'],
        slug := <str>data['slug'],
        description := <str>json_get(data, 'description'),
				maintainers := (
					select people::Person
					filter .alias in <str>json_array_unpack(data['maintainers'])
				) ?? (SELECT admin::Settings.superadmin.identity),
        occurrences := (
          select occurrence::Occurrence
          filter .id in <uuid>json_array_unpack(occurrences)['id']
        )
      }) {
        *,
				maintainers: { * },
				occurrences: {
					id, comments,
					sampling: {
						*,
						target_taxa: { * },
						methods: { * },
						fixatives: { * },
						habitats: { * },
					},
					identification: { ** },
					category:= [is occurrence::BioMaterial].category ?? occurrence::OccurrenceCategory.External,
					element:= (
						if exists [is occurrence::BioMaterial].id
						then "BioMaterial"
						else "Sequence"
					),

				}
      }
      `, &created, data, occurrencesData,
	)
	return
}
