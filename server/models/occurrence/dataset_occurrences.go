package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/dataset"
)

type OccurrenceDataset struct {
	dataset.Dataset `gel:"$inline" json:",inline"`
	Sites           []SiteItem               `gel:"sites" json:"sites"`
	Occurrences     []OccurrenceWithCategory `gel:"occurrences" json:"occurrences"`
	IsCongruent     bool                     `gel:"is_congruent" json:"is_congruent"`
}

func ListOccurrenceDatasets(db geltypes.Executor) ([]OccurrenceDataset, error) {
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

func GetOccurrenceDataset(db geltypes.Executor, slug string) (dataset OccurrenceDataset, err error) {
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
	dataset.DatasetInput `gel:"$inline" json:",inline"`
	Occurrences          OccurrenceBatchInput `json:"occurrences"`
}

func (i OccurrenceDatasetInput) SaveTx(tx geltypes.Tx) (created OccurrenceDataset, err error) {
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
