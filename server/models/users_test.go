package models

import (
	"github.com/lsdch/biome/db/biomedb"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUserRoleAuthorization(t *testing.T) {
	for _, tt := range []struct {
		role, requested UserRole
		allowed         bool
	}{
		{biomedb.UserRoleVisitor, biomedb.UserRoleContributor, false},
		{biomedb.UserRoleContributor, biomedb.UserRoleContributor, true},
		{biomedb.UserRoleMaintainer, biomedb.UserRoleAdmin, false},
		{biomedb.UserRoleAdmin, biomedb.UserRoleMaintainer, true},
	} {
		t.Run(string(tt.role)+"/"+string(tt.requested), func(t *testing.T) {
			user := User{Role: tt.role}
			require.Equal(t, tt.allowed, user.IsGranted(tt.requested))
		})
	}
}

func TestInvitationRegistrationHashesPassword(t *testing.T) {
	input := RegisterFromInvitationParams{Login: "alice", Password: "test-password-42"}
	params, err := input.ToParams("invitation-hash")
	require.NoError(t, err)
	require.Equal(t, "invitation-hash", params.TokenHash)
	require.Equal(t, input.Login, params.Login)
	require.NotEqual(t, input.Password, params.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(params.PasswordHash), []byte(input.Password)))
}
