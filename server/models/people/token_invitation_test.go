package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/tests"
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

func TestValidateInvitation(t *testing.T) {
	client := db.Client()
	person := SetupPerson(t, client)
	invitation, err := person.CreateInvitation(client, people.Contributor)
	require.NoError(t, err)
	i, err := people.ValidateInvitationToken(client, invitation.Token)
	require.NoError(t, err)
	assert.Equal(t, person.FullName, i.Person.FullName)
}

func TestClaimInvitation(t *testing.T) {
	client := db.Client()
	input := tests.FakeData[people.UserInput](t)
	person := SetupPerson(t, client)
	invitation, err := person.CreateInvitation(client, people.Contributor)
	require.NoError(t, err)
	u, err := input.ClaimInvitationToken(client, invitation.Token)
	require.NoError(t, err)
	assert.Equal(t, person.ID, u.Person.ID)
	assert.Equal(t, invitation.Role, u.Role)
}
