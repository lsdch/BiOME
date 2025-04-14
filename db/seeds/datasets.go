package seeds

import (
	"encoding/json"
	"fmt"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/occurrence"

	"github.com/sirupsen/logrus"
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

func LoadMultipleOccurrencesDatasets(file string) (datasets []occurrence.OccurrenceDatasetInput, err error) {
	b, err := data.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &datasets)
	return
}

func SeedOccurrencesDatasets(tx geltypes.Tx, datasets []occurrence.OccurrenceDatasetInput) error {

	for _, dataset := range datasets {
		created, err := dataset.SaveTx(tx)
		if err != nil {
			return fmt.Errorf("❗Failed to seed occurrence dataset: %v", err)
		}
		logrus.Infof("🌱 dataset: %s", created.Label)
	}
	return nil
}
