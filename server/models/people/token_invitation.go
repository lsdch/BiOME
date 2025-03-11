package people

import (
	"context"
	"net/url"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/settings"
	"github.com/lsdch/biome/models/tokens"
	"github.com/lsdch/biome/services/email"
	email_templates "github.com/lsdch/biome/templates"
)

type InvitationOptions struct {
	Email string   `gel:"email" json:"email" format:"email" doc:"E-mail address of the recipient of the invitation"`
	Role  UserRole `gel:"role" json:"role"`
}

type InvitationInput struct {
	Person             PersonInner `gel:"identity" json:"identity"`
	InvitationOptions  `gel:"$inline" json:",inline"`
	tokens.TokenRecord `gel:"$inline" json:",inline"`
}

type Invitation struct {
	ID              geltypes.UUID `gel:"id" json:"id"`
	IssuedBy        User          `gel:"issued_by" json:"issued_by"`
	InvitationInput `gel:"$inline" json:",inline"`
}

func (i InvitationInput) Save(db geltypes.Executor) (Invitation, error) {
	var invitation Invitation
	err := db.QuerySingle(context.Background(),
		`#edgeql
			with module people
			select (insert tokens::UserInvitation {
				identity := (select Person filter .id = <uuid>$0),
				role := <UserRole>$1,
				token := <str>$2,
				expires := <datetime>$3,
				email := <str>$4
			}) { ** }
		`,
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

func ValidateInvitationToken(db geltypes.Executor, token tokens.Token) (Invitation, error) {
	var invitation Invitation
	err := db.QuerySingle(context.Background(),
		`#edgeql
			select tokens::UserInvitation { ** } filter .token = <str>$0
		`, &invitation, token,
	)
	return invitation, err
}
