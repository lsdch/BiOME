package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/tests"
	"testing"

	"github.com/edgedb/edgedb-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupAuthAdminUser(t *testing.T) edgedb.Executor {
	client := db.Client()
	user := SetupUser(t)
	require.NoError(t, user.SetRole(client, people.Admin))
	auth_client := client.WithGlobals(map[string]interface{}{
		"current_user_id": user.ID,
	})
	return auth_client
}

func TestInvitation(t *testing.T) {
	client := db.Client()
	person := SetupPerson(t, client)
	auth_client := SetupAuthAdminUser(t)
	invitation, err := person.CreateInvitation(
		people.InvitationOptions{Role: people.Maintainer},
	).Save(auth_client)
	require.NoError(t, err)
	assert.True(t, invitation.IsValid())
}

func TestValidateInvitation(t *testing.T) {
	client := db.Client()
	person := SetupPerson(t, client)

	auth_client := SetupAuthAdminUser(t)
	invitation, err := person.CreateInvitation(
		people.InvitationOptions{Role: people.Maintainer},
	).Save(auth_client)
	require.NoError(t, err)
	i, err := people.ValidateInvitationToken(client, invitation.Token)
	require.NoError(t, err)
	assert.Equal(t, person.FullName, i.Person.FullName)
}

func TestClaimInvitation(t *testing.T) {
	client := db.Client()
	input := tests.FakeData[people.UserInput](t)
	person := SetupPerson(t, client)

	auth_client := SetupAuthAdminUser(t)
	invitation, err := person.CreateInvitation(
		people.InvitationOptions{Role: people.Maintainer},
	).Save(auth_client)
	require.NoError(t, err)
	u, err := input.ClaimInvitationToken(client, invitation.Token)
	require.NoError(t, err)
	assert.Equal(t, person.ID, u.Person.ID)
	assert.Equal(t, invitation.Role, u.Role)
}
