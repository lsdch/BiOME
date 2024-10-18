package people

import (
	"context"
	"darco/proto/models/settings"
	"darco/proto/models/tokens"
	"darco/proto/services/email"
	email_templates "darco/proto/templates"
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
		select (insert tokens::UserInvitation {
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
	token := tokens.GenerateToken(settings.Security().InvitationTokenDuration())

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

	templateData := email_templates.InvitationData{
		Name:       i.Person.FirstName,
		IssuerName: i.IssuedBy.Person.FullName,
		App:        settings.Instance().Name,
		Role:       string(i.Role),
		URL:        target,
	}

	sendError := (&email.EmailData{
		To:       i.Email,
		From:     settings.Email().FromHeader(),
		Subject:  templateData.Subject(),
		Template: email_templates.Invitation(templateData),
	}).Send(settings.Email().FromHeader())
	return &target, sendError
}

func ValidateInvitationToken(db edgedb.Executor, token tokens.Token) (Invitation, error) {
	var invitation Invitation
	err := db.QuerySingle(context.Background(),
		`select tokens::UserInvitation { ** } filter .token = <str>$0`,
		&invitation, token,
	)
	return invitation, err
}
