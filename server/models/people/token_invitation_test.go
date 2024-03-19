package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvitation(t *testing.T) {
	client := db.Client()
	person := SetupPerson(t, client)
	invitation, err := person.CreateInvitation(client, people.Maintainer)
	require.NoError(t, err)
	assert.True(t, invitation.IsValid())
}
