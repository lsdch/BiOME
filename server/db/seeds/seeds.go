package main

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/taxonomy"
	"fmt"
	"os"

	"github.com/edgedb/edgedb-go"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

func entityQueryPath(entity string) string {
	return fmt.Sprintf("queries/%s.edgeql", entity)
}
func entityDataPath(entity string) string {
	return fmt.Sprintf("data/%s.json", entity)
}

func Seed(tx *edgedb.Tx, entity string) error {
	queryPath := entityQueryPath(entity)
	dataPath := entityDataPath(entity)
	query, err := os.ReadFile(queryPath)
	if err != nil {
		logrus.Errorf("Failed to load seed query @ %s: %v", queryPath, err)
		return err
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		logrus.Errorf("Failed to load seed data @ %s: %v", dataPath, err)
		return err
	}

	err = tx.Execute(context.Background(), string(query), data)
	if err != nil {
		logrus.Errorf(
			"Query execution failed for query @ %s and data @ %s:\n%v",
			queryPath, dataPath, err,
		)
		return err
	}
	return nil
}

func seedTaxonomyGBIF() error {
	bar := progressbar.Default(-1, "Importing Asellidae taxonomy from GBIF")
	models.DB().Execute(context.Background(), "delete taxonomy::Taxon")
	var total int
	err := taxonomy.ImportTaxon(models.DB(), 4574, func(p *taxonomy.ImportProcess) {
		total = p.Imported
		bar.Set(p.Imported)
	})
	bar.Clear()

	logrus.Infof("Taxonomy setup done: %d taxa imported", total)
	return err
}

var entities = []string{
	"institutions",
	"persons",
}

func main() {
	if err := seedTaxonomyGBIF(); err != nil {
		logrus.Errorf("Failed to load Asellidae taxonomy: %v", err)
		return
	}

	models.DB().Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
		for _, entity := range entities {
			logrus.Infof("Seeding %s", entity)
			err := Seed(tx, entity)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
