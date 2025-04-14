package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/dataset"
)

type OccurrenceDatasetListItem struct {
	dataset.Dataset `gel:"$inline" json:",inline"`
	Sites           int64 `gel:"sites_count" json:"sites"`
	Occurrences     int64 `gel:"occurrences_count" json:"occurrences"`
	IsCongruent     bool  `gel:"is_congruent" json:"is_congruent"`
}

func ListOccurrenceDatasets(db geltypes.Executor) ([]OccurrenceDatasetListItem, error) {
	datasets := []OccurrenceDatasetListItem{}
	err := db.Query(context.Background(),
		`#edgeql
			select datasets::OccurrenceDataset {
				*,
				maintainers: { *, user: { * } },
				meta: { * },
				sites_count := count(.sites),
				occurrences_count := count(.occurrences),
			}
		`,
		&datasets,
	)
	return datasets, err
}

type OccurrenceDataset struct {
	dataset.Dataset `gel:"$inline" json:",inline"`
	Sites           []SiteItem               `gel:"sites" json:"sites"`
	Occurrences     []OccurrenceWithCategory `gel:"occurrences" json:"occurrences"`
	IsCongruent     bool                     `gel:"is_congruent" json:"is_congruent"`
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
	OccurrenceBatchInput `json:",inline"`
}

func (i OccurrenceDatasetInput) SaveTx(tx geltypes.Tx) (created OccurrenceDataset, err error) {
	occurrences, err := i.OccurrenceBatchInput.Save(tx)
	if err != nil {
		return created, err
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
