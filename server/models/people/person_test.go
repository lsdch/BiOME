package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/tests"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed fixtures/person.json
var mockPersonJSON string

func FakePersonInput(t *testing.T) *people.PersonInput {
	return tests.FakeData[people.PersonInput](t)
}

func TestPerson(t *testing.T) {
	t.Run("Create person", func(t *testing.T) {
		_, err := FakePersonInput(t).Create(db.Client())
		require.NoError(t, err)
	})
}
