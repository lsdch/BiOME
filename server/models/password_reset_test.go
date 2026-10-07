package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordResetParamsValidate(t *testing.T) {
	user := User{Login: "alice", Email: "alice@example.org", FirstName: "Alice", LastName: "Test"}
	for _, password := range []string{"", "password", strings.Repeat("a", 73), strings.Repeat("é", 37)} {
		require.Error(t, (PasswordResetParams{NewPassword: password}).Validate(user, 3))
	}
	params := PasswordResetParams{NewPassword: "Lunar!Cactus-87.River$Quartz"}
	require.NoError(t, params.Validate(user, 3))
}
