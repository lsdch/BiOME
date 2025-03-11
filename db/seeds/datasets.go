package seeds

import (
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/occurrence"

	"github.com/sirupsen/logrus"
)

func LoadOccurrencesDatasets(file string) (datasets []occurrence.OccurrenceDatasetInput, err error) {
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
			return err
		}
		logrus.Infof("🌱 dataset: %s", created.Label)
	}
	return nil
}
