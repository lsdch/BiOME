package seeds

import (
	"encoding/json"

	"github.com/lsdch/biome/models/occurrence"
)

func LoadOccurrencesDataset(file string) (*occurrence.OccurrenceDatasetInput, error) {
	dataset := new(occurrence.OccurrenceDatasetInput)
	b, err := data.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, dataset)
	return dataset, err
}

func LoadMultipleOccurrencesDatasets(file string) (datasets []*occurrence.OccurrenceDatasetInput, err error) {
	b, err := data.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &datasets)
	return
}
