package people

import (
	"context"
	"darco/proto/models/settings"
	"darco/proto/models/tokens"
	"darco/proto/services/email"
	"fmt"
	"net/url"

	"github.com/edgedb/edgedb-go"
)

type InvitationOptions struct {
	Email string   `edgedb:"email" json:"email" format:"email" doc:"E-mail address of the recipient of the invitation"`
	Role  UserRole `edgedb:"role" json:"role"`
}

type InvitationInput struct {
	Person             PersonInner `edgedb:"identity" json:"identity"`
	InvitationOptions  `edgedb:"$inline" json:",inline"`
	tokens.TokenRecord `edgedb:"$inline" json:",inline"`
}

type Invitation struct {
	ID              edgedb.UUID `edgedb:"id" json:"id"`
	IssuedBy        User        `edgedb:"issued_by" json:"issued_by"`
	InvitationInput `edgedb:"$inline" json:",inline"`
}

func (i InvitationInput) Save(db edgedb.Executor) (Invitation, error) {
	var invitation Invitation
	err := db.QuerySingle(context.Background(),
		`with module people
		select (insert UserInvitation {
			identity := (select Person filter .id = <uuid>$0),
			role := <UserRole>$1,
			token := <str>$2,
			expires := <datetime>$3,
			email := <str>$4
		}) { ** }`,
		&invitation,
		i.Person.ID, i.Role, i.Token, i.Expires, i.Email,
	)
	return invitation, err
}

func (p *Person) CreateInvitation(options InvitationOptions) InvitationInput {
	token := tokens.GenerateToken(20)

	return InvitationInput{
		Person:            p.PersonInner,
		InvitationOptions: options,
		TokenRecord:       token,
	}
}

// Sends the invitation token to the assigned email address,
// with an activation link derived from the `target` argument.
func (i *Invitation) Send(target url.URL) (*url.URL, error) {
	params := target.Query()
	params.Set("token", string(i.Token))
	target.RawQuery = params.Encode()

	sendError := (&email.EmailData{
		To:       i.Email,
		Subject:  fmt.Sprintf("Invitation to create an account on %s", settings.Instance().Name),
		Template: "invitation.html",
		Data: map[string]any{
			"Name":       i.Person.FullName,
			"IssuerName": i.IssuedBy.Person.FullName,
			"App":        settings.Instance().Name,
			"Role":       i.Role,
			"URL":        target.String(),
		},
	}).Send(email.AdminEmailAddress())
	return &target, sendError
}

func ValidateInvitationToken(db edgedb.Executor, token tokens.Token) (Invitation, error) {
	var invitation Invitation
	err := db.QuerySingle(context.Background(),
		`select people::UserInvitation { ** } filter .token = <str>$0`,
		&invitation, token,
	)
	return invitation, err
}
