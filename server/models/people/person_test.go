package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/tests"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/edgedb/edgedb-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed fixtures/person.json
var mockPersonJSON string

func FakePersonInput(t *testing.T) *people.PersonInput {
	return tests.FakeData[people.PersonInput](t)
}

func SetupPerson(t *testing.T, db edgedb.Executor) people.Person {
	p, err := FakePersonInput(t).Create(db)
	require.NoError(t, err)
	return p
}

func TestPerson(t *testing.T) {
	client := db.Client()
	t.Run("Create person", func(t *testing.T) {
		t.Parallel()
		input := FakePersonInput(t)
		alias := input.GenerateAlias()
		p, err := input.Create(client)
		require.NoError(t, err)
		assert.Equal(t, p.Alias, *alias)
	})

	t.Run("Delete person", func(t *testing.T) {
		t.Parallel()
		p := SetupPerson(t, client)
		deleted, err := p.Delete(client)
		require.NoError(t, err)
		assert.Equal(t, p, deleted)
	})

	t.Run("Find person", func(t *testing.T) {
		t.Parallel()
		p := SetupPerson(t, client)
		found, err := people.FindPerson(client, p.ID)
		require.NoError(t, err)
		assert.Equal(t, p, found)
	})

	t.Run("List persons", func(t *testing.T) {
		t.Parallel()
		_ = SetupPerson(t, client)
		persons, err := people.ListPersons(client)
		require.NoError(t, err)
		assert.NotEmpty(t, persons)
	})

	t.Run("Update person", func(t *testing.T) {
		t.Parallel()
		p := SetupPerson(t, client)
		u := tests.FakeData[people.PersonUpdate](t)
		json, _ := json.Marshal(u)
		id, err := u.Update(client, p.ID)
		require.NoErrorf(t, err, "%s", json)
		assert.Equal(t, p.ID, id)
		p, err = people.FindPerson(client, p.ID)
		require.NoError(t, err)
		assert.Equal(t, p.FirstName, *u.FirstName)
	})
}
