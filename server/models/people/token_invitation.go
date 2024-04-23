package people

import (
	"context"
	"darco/proto/models/settings"
	"darco/proto/services/email"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type InvitationOptions struct {
	Dest string   `edgedb:"dest" json:"dest" format:"email"`
	Role UserRole `edgedb:"role" json:"role"`
}

type InvitationInput struct {
	Person            PersonInner `edgedb:"identity" json:"identity"`
	InvitationOptions `edgedb:"$inline" json:",inline"`
	TokenWrapper      `edgedb:"$inline" json:",inline"`
}

type Invitation struct {
	ID              edgedb.UUID `edgedb:"id" json:"id"`
	IssuedBy        User        `edgedb:"issued_by" json:"issued_by"`
	InvitationInput `edgedb:"$inline" json:",inline"`
}

func (i InvitationInput) Save(db edgedb.Executor) (*Invitation, error) {
	var invitation Invitation
	if err := db.QuerySingle(context.Background(),
		`with module people
		select (insert UserInvitation {
			identity := (select Person filter .id = <uuid>$0),
			role := <UserRole>$1,
			token := <str>$2,
			expires := <datetime>$3,
			dest := <str>$4
		}) { ** }`, &invitation,
		i.Person.ID, i.Role, i.Token, i.Expires, i.Dest,
	); err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (p *Person) CreateInvitation(options InvitationOptions) InvitationInput {
	tokStr := GenerateToken(20)
	expires := time.Now().Add(settings.Security().AccountTokenDuration())

	return InvitationInput{
		Person:            p.PersonInner,
		InvitationOptions: options,
		TokenWrapper: TokenWrapper{
			Token:   tokStr,
			Expires: expires,
		},
	}
}

const InvitationTokenPathParameter = "{token}"

var tokenParamRegexCheck *regexp.Regexp

func init() {
	re, err := regexp.Compile(`\{.*\}`)
	if err != nil {
		logrus.Fatalf("Failed to compile token URL verification regex: %v", err)
	}
	tokenParamRegexCheck = re
}

func (i *Invitation) encodeTokenURL(target url.URL) (*url.URL, error) {
	if !strings.Contains(target.Path, InvitationTokenPathParameter) {
		return nil, fmt.Errorf("Target URL path to activate invitation token lacks '%s' parameter.", InvitationTokenPathParameter)
	}
	target.Path = strings.Replace(
		target.Path,
		InvitationTokenPathParameter,
		string(i.Token),
		1,
	)
	paramsMissing := tokenParamRegexCheck.Match([]byte(target.Path))
	if paramsMissing {
		return nil, fmt.Errorf("Some parameters in the token URL path were not provided: %v", target.Path)
	}
	return &target, nil
}

// Sends the invitation token to the assigned email address,
// with an activation link derived from the `target` argument.
// Target URL must include a `{token}` path argument.
func (i *Invitation) Send(target url.URL) (*url.URL, error) {
	activationURL, err := i.encodeTokenURL(target)
	if err != nil {
		return nil, err
	}
	data := map[string]any{
		"Name":       i.Person.FullName,
		"IssuerName": i.IssuedBy.Person.FullName,
		"App":        settings.Instance().Name,
		"Role":       i.Role,
		"URL":        activationURL.String(),
	}
	return activationURL, email.Send(email.AdminEmailAddress(),
		&email.EmailData{
			To:       i.Dest,
			Subject:  fmt.Sprintf("Invitation to collaborate on %s", settings.Instance().Name),
			Template: "invitation.html",
			Data:     data,
		})
}

func ValidateInvitationToken(db edgedb.Executor, token Token) (*Invitation, error) {
	var invitation Invitation
	err := db.QuerySingle(context.Background(),
		`select people::UserInvitation { ** } filter .token = <str>$0`,
		&invitation, token,
	)
	return &invitation, err
}

var InvalidTokenError = fmt.Errorf("Invalid token")

func (u *UserInput) ClaimInvitationToken(db edgedb.Executor, token Token) (*User, error) {
	invitation, err := ValidateInvitationToken(db, token)
	if err != nil {
		return nil, InvalidTokenError
	}
	user, err := u.Save(db, invitation.Role)
	if err != nil {
		return nil, fmt.Errorf("User registration failed: %w", err)
	}
	if err := user.SetIdentity(db, &invitation.Person); err != nil {
		return nil, fmt.Errorf("Failed to assign identity to registered user: %w", err)
	}
	return user, nil
}
