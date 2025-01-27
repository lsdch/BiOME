package people_test

import (
	"testing"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeta(t *testing.T) {
	user, err := people.Find(db.Client(), "dev.admin")
	require.NoError(t, err)
	client := db.WithCurrentUser(user.ID)
	t.Run("Meta update", func(t *testing.T) {
		t.Parallel()
		person := SetupPerson(t, client)
		updateInput := tests.FakeData[people.PersonUpdate](t)
		_, err := updateInput.Save(client, person.ID)
		require.NoError(t, err)
		updated, err := people.FindPerson(client, person.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, updated.Meta.UpdatedBy.ID)
		assert.GreaterOrEqual(t, updated.Meta.LastUpdated, updated.Meta.Created)
	})
}
