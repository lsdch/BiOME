package main

import (
	"context"
	gbif "darco/proto/models/taxonomy/GBIF"
	"fmt"

	"github.com/edgedb/edgedb-go"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

func seedTaxonomyGBIF(db *edgedb.Client) error {
	bar := progressbar.Default(-1, "Importing Asellidae taxonomy from GBIF")
	db.Execute(context.Background(), "delete taxonomy::Taxon")
	var total int
	err := gbif.ImportTaxon(db,
		gbif.ImportRequestGBIF{Key: 4574, Children: true},
		func(p *gbif.ImportProcess) {
			total = p.Imported
			bar.Set(p.Imported)
			if p.Error != nil {
				bar.Describe(fmt.Sprintf("%+v", p))
				logrus.Fatalf("Failed to import taxonomy: %v", p.Error)
			}
		})
	bar.Clear()

	logrus.Infof("Taxonomy setup done: %d taxa imported", total)
	return err
}
