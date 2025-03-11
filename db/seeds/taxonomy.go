package seeds

import (
	"context"
	"fmt"

	gbif "github.com/lsdch/biome/models/taxonomy/GBIF"

	"github.com/geldata/gel-go/geltypes"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

func SeedTaxonomyGBIF(db geltypes.Tx) error {
	bar := progressbar.Default(-1, "Importing Asellidae taxonomy from GBIF")
	_ = db.Execute(context.Background(), "delete taxonomy::Taxon")
	var total int
	err := gbif.ImportTaxonTx(db,
		gbif.ImportRequestGBIF{Key: 4574, Children: true},
		func(p *gbif.ImportProcess) {
			total = p.Imported
			_ = bar.Set(p.Imported)
			if p.Error != nil {
				bar.Describe(fmt.Sprintf("%+v", p))
				logrus.Fatalf("Failed to import taxonomy: %v", p.Error)
			}
		})
	_ = bar.Close()

	logrus.Infof("Taxonomy setup done: %d taxa imported", total)
	return err
}
