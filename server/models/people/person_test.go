package people_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FakePersonInput(t *testing.T) *people.PersonInput {
	return tests.FakeData[people.PersonInput](t)
}

func SetupPerson(t *testing.T, db geltypes.Executor) people.Person {
	p, err := FakePersonInput(t).Save(db)
	require.NoError(t, err)
	return p
}

func TestPerson(t *testing.T) {
	client := db.Client()
	t.Run("Create person", func(t *testing.T) {
		input := FakePersonInput(t)
		alias := input.GenerateAlias()
		p, err := input.Save(client)
		require.NoError(t, err)
		assert.Equal(t, p.Alias, alias)
	})

	t.Run("Delete person", func(t *testing.T) {
		p := SetupPerson(t, client)
		deleted, err := p.Delete(client)
		require.NoError(t, err)
		assert.Equal(t, p, deleted)
	})

	t.Run("Find person", func(t *testing.T) {
		p := SetupPerson(t, client)
		found, err := people.FindPerson(client, p.ID)
		require.NoError(t, err)
		assert.Equal(t, p, found)
	})

	t.Run("List persons", func(t *testing.T) {
		_ = SetupPerson(t, client)
		persons, err := people.ListPersons(client)
		require.NoError(t, err)
		assert.NotEmpty(t, persons)
	})

	t.Run("Update person", func(t *testing.T) {
		p := SetupPerson(t, client)
		u := tests.FakeData[people.PersonUpdate](t)
		u.FirstName.IsSet = true
		json, _ := json.Marshal(u)
		updated, err := u.Save(client, p.ID)
		require.NoErrorf(t, err, "%s", json)
		assert.Equal(t, p.ID, updated.ID)
		assert.Equal(t, u.FirstName.Value, updated.FirstName)
	})
}
