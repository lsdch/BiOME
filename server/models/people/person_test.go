package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/tests"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed fixtures/person.json
var mockPersonJSON string

func FakePersonInput(t *testing.T) *people.PersonInput {
	return tests.FakeData[people.PersonInput](t)
}

func TestPerson(t *testing.T) {
	client := db.Client()
	t.Run("Create person", func(t *testing.T) {
		_, err := FakePersonInput(t).Create(client)
		require.NoError(t, err)
	})

	t.Run("Delete person", func(t *testing.T) {
		p, err := FakePersonInput(t).Create(client)
		require.NoError(t, err)
		deleted, err := p.Delete(client)
		require.NoError(t, err)
		assert.Equal(t, p, deleted)

	})
}
