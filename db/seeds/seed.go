package seeds

import (
	"context"
	"embed"
	"fmt"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

//go:embed queries
var queries embed.FS

//go:embed data
var data embed.FS

func entityQueryPath(entity string) string {
	return fmt.Sprintf("queries/%s.edgeql", entity)
}
func entityDataPath(entity string) string {
	return fmt.Sprintf("data/%s.json", entity)
}

func Seed(tx *edgedb.Tx, entity string) error {
	queryPath := entityQueryPath(entity)
	dataPath := entityDataPath(entity)
	query, err := queries.ReadFile(queryPath)
	if err != nil {
		logrus.Errorf("Failed to load seed query @ %s: %v", queryPath, err)
		return err
	}

	data, err := data.ReadFile(dataPath)
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
