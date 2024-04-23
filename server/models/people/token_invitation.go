package people

import (
	"context"
	"darco/proto/models/settings"
	"time"

	"github.com/edgedb/edgedb-go"
)

type InvitationOptions struct {
	Address string   `edgedb:"address" json:"address" format:"email"`
	Role    UserRole `edgedb:"role" json:"role"`
}

type InvitationInput struct {
	Person       PersonInner `edgedb:"identity" json:"identity"`
	Role         UserRole    `edgedb:"role" json:"role"`
	TokenWrapper `edgedb:"$inline" json:",inline"`
} // @name InvitationInput

type Invitation struct {
	ID              edgedb.UUID `edgedb:"id" json:"id"`
	IssuedBy        User        `edgedb:"issued_by" json:"issued_by"`
	InvitationInput `edgedb:"$inline" json:",inline"`
} // @name Invitation

func (i InvitationInput) Save(db edgedb.Executor) (*Invitation, error) {
	var invitation Invitation
	err := db.QuerySingle(context.Background(),
		`with module people
		select (insert UserInvitation {
			identity := (select Person filter .id = <uuid>$0),
			role := <UserRole>$1,
			token := <str>$2,
			expires := <datetime>$3
		}) { ** }`, &invitation,
		i.Person.ID, i.Role, i.Token, i.Expires,
	)
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (p *Person) CreateInvitation(db edgedb.Executor, role UserRole) (*Invitation, error) {
	tokStr := GenerateToken(20)
	expires := time.Now().Add(settings.Security().AccountTokenDuration())
	token := InvitationInput{
		Person: p.PersonInner,
		Role:   role,
		TokenWrapper: TokenWrapper{
			Token:   tokStr,
			Expires: expires,
		},
	}
	return token.Save(db)
}

func ValidateInvitationToken(db edgedb.Executor, token Token) (*Invitation, error) {
	var invitation Invitation
	err := db.QuerySingle(context.Background(),
		`select people::UserInvitation { ** } filter .token = <str>$0`,
		&invitation, token,
	)
	return &invitation, err
}

type InvitationClaimErrorType string

const (
	InvalidToken          InvitationClaimErrorType = "Invalid token"
	RegistrationFailed    InvitationClaimErrorType = "User registration failed"
	IdentityLinkageFailed InvitationClaimErrorType = "Failed to assign identity to registered user "
)

type InvitationClaimError struct {
	Type InvitationClaimErrorType
	Err  error
}

func (err InvitationClaimError) Error() string {
	return err.Err.Error()
}

func (u *UserInput) ClaimInvitationToken(db edgedb.Executor, token Token) (*User, *InvitationClaimError) {
	invitation, err := ValidateInvitationToken(db, token)
	if err != nil {
		return nil, &InvitationClaimError{Type: InvalidToken, Err: err}
	}
	user, err := u.Save(db, invitation.Role)
	if err != nil {
		return nil, &InvitationClaimError{Type: RegistrationFailed, Err: err}
	}
	if err := user.SetIdentity(db, &invitation.Person); err != nil {
		return nil, &InvitationClaimError{Type: IdentityLinkageFailed, Err: err}
	}
	return user, nil
}
