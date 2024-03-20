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
