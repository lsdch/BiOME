package models

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestRegisterFromInvitationValidation(t *testing.T) {
	const strongPassword = "r8!Fz2#Lq9$Vn4@Tk7%W"
	valid := RegisterFromInvitationParams{
		Login: "alice", FirstName: "Alice", LastName: "Example", Password: strongPassword,
	}
	for _, tt := range []struct {
		name       string
		change     func(*RegisterFromInvitationParams)
		field, tag string
	}{
		{name: "valid", change: func(p *RegisterFromInvitationParams) {}},
		{name: "unicode names", change: func(p *RegisterFromInvitationParams) { p.FirstName = "李雷"; p.LastName = "杜甫" }},
		{name: "surrounding whitespace", change: func(p *RegisterFromInvitationParams) { p.Login = " alice "; p.FirstName = " Alice " }},
		{name: "empty login", change: func(p *RegisterFromInvitationParams) { p.Login = "" }, field: "Login", tag: "required"},
		{name: "blank first name", change: func(p *RegisterFromInvitationParams) { p.FirstName = " \t\n" }, field: "FirstName", tag: "required"},
		{name: "short trimmed login", change: func(p *RegisterFromInvitationParams) { p.Login = " a " }, field: "Login", tag: "min"},
		{name: "short unicode name", change: func(p *RegisterFromInvitationParams) { p.LastName = " é " }, field: "LastName", tag: "min"},
		{name: "empty password", change: func(p *RegisterFromInvitationParams) { p.Password = "" }, field: "Password", tag: "required"},
		{name: "weak password", change: func(p *RegisterFromInvitationParams) { p.Password = "password" }, field: "Password", tag: "password_strength"},
		{name: "72 bytes", change: func(p *RegisterFromInvitationParams) { p.Password += strings.Repeat("x", 72-len(strongPassword)) }},
		{name: "73 bytes", change: func(p *RegisterFromInvitationParams) { p.Password += strings.Repeat("x", 73-len(strongPassword)) }, field: "Password", tag: "maxbytes"},
		{name: "unicode bytes", change: func(p *RegisterFromInvitationParams) { p.Password = strings.Repeat("é", 37) }, field: "Password", tag: "maxbytes"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input := valid
			tt.change(&input)
			before := input
			err := input.Validate("alice@example.org", 3)
			require.Equal(t, before, input, "validation must preserve supplied values")
			if tt.field == "" {
				require.NoError(t, err)
				return
			}
			var validationErrors validator.ValidationErrors
			require.ErrorAs(t, err, &validationErrors)
			require.Len(t, validationErrors, 1)
			require.Equal(t, tt.field, validationErrors[0].Field())
			require.Equal(t, tt.tag, validationErrors[0].Tag())
		})
	}
}

func TestRegisterFromInvitationPasswordUsesInvitedEmail(t *testing.T) {
	input := RegisterFromInvitationParams{
		Login: "alice", FirstName: "Alice", LastName: "Example", Password: "r8!Fz2#Lq9$Vn4@Tk7%W",
	}
	require.NoError(t, input.Validate("alice@example.org", 3))
	// The invited email is a personal-information dictionary entry, not caller-supplied data.
	err := input.Validate(input.Password, 3)
	var validationErrors validator.ValidationErrors
	require.ErrorAs(t, err, &validationErrors)
	require.Equal(t, "password_strength", validationErrors[0].Tag())
}
