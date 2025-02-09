package people

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/settings"
	"github.com/lsdch/biome/models/tokens"
	"github.com/lsdch/biome/services/email"
	email_templates "github.com/lsdch/biome/templates"

	"github.com/edgedb/edgedb-go"
)

type EmailField struct {
	Email string `edgedb:"email" json:"email" format:"email" fake:"{email}"`
}

type UserInput struct {
	Login         string `edgedb:"login" json:"login" binding:"login,required,unique_login" fake:"{username}"`
	EmailField    `edgedb:"$inline" json:",inline"`
	PasswordInput `json:",inline"`
}

func (u UserInput) Save(db edgedb.Executor, role UserRole, identity PersonInner) (*User, error) {
	var user User
	input, _ := json.Marshal(u)
	err := db.QuerySingle(context.Background(),
		`#edgeql
			with module people,
			data := <json>$0,
			user := (insert User {
				login := <str>data['login'],
				email := <str>data['email'],
				password := <str>data['password'],
				role := <UserRole>$1,
				identity := assert_single((select Person filter .id = <uuid>$2))
			}),
			select user { ** }
		`,
		&user, input, role, identity.ID,
	)
	return &user, err
}

var InvalidTokenError = fmt.Errorf("Invalid token")

func (u UserInput) RegisterWithToken(db edgedb.Executor, token tokens.Token) (*User, error) {
	invitation, err := ValidateInvitationToken(db, token)
	if err != nil {
		return nil, InvalidTokenError
	}
	user, err := u.Save(db, invitation.Role, invitation.Person)
	if err != nil {
		return nil, fmt.Errorf("User registration failed: %w", err)
	}
	return user, nil
}

type PendingUserRequestInput struct {
	EmailField     `json:",inline" edgedb:"$inline"`
	PersonIdentity `edgedb:"$inline" json:",inline"`
	Organisation   string `json:"organisation,omitempty" edgedb:"organisation" fake:"{word}"`
	Motive         string `json:"motive,omitempty" edgedb:"motive" fake:"{sentence:10}"`
}

//go:embed queries/register_pending_user.edgeql
var registerPendingUserQuery string

// Creates a request for a user account which can be validated by and admin
// to send an invitation to create an account
func (u *PendingUserRequestInput) Register(db edgedb.Executor) (*PendingUserRequest, error) {
	args, _ := json.Marshal(u)
	var pendingUser PendingUserRequest
	err := db.QuerySingle(context.Background(), registerPendingUserQuery, &pendingUser, args)
	return &pendingUser, err
}

type PendingUserRequest struct {
	ID             edgedb.UUID `edgedb:"id" json:"id"`
	EmailField     `json:",inline" edgedb:"$inline"`
	PersonIdentity `edgedb:"$inline" json:",inline"`
	FullName       string             `edgedb:"full_name" json:"full_name"`
	Organisation   edgedb.OptionalStr `json:"organisation,omitempty" edgedb:"organisation"`
	Motive         edgedb.OptionalStr `json:"motive,omitempty" edgedb:"motive"`
	CreatedOn      time.Time          `json:"created_on" edgedb:"created_on"`
	EmailVerified  bool               `edgedb:"email_verified" json:"email_verified"`
}

func (p *PendingUserRequest) Delete(db edgedb.Executor) error {
	return db.Execute(context.Background(),
		`#edgeql
			delete <people::PendingUserRequest><uuid>$0;
		`, p.ID,
	)
}

func (p *PendingUserRequest) SetEmailVerified(db edgedb.Executor, isVerified bool) error {
	err := db.Execute(context.Background(),
		`#edgeql
			update <people::PendingUserRequest><uuid>$0 set { email_verified := <bool>$1 }
		`, p.ID, isVerified,
	)
	if err != nil {
		return err
	}
	p.EmailVerified = true
	return nil
}

func ListPendingUserRequests(db edgedb.Executor) ([]PendingUserRequest, error) {
	var items = []PendingUserRequest{}
	err := db.Query(context.Background(),
		`#edgeql
			select people::PendingUserRequest { ** } order by .created_on desc;
		`, &items,
	)
	return items, err
}

func GetPendingUserRequest(db edgedb.Executor, email string) (*PendingUserRequest, error) {
	var req PendingUserRequest
	err := db.QuerySingle(context.Background(),
		`#edgeql
			select people::PendingUserRequest { ** } filter .email = <str>$0;
		`, &req, email,
	)
	return &req, err
}

func DeletePendingUserRequest(db edgedb.Executor, email string) (deleted PendingUserRequest, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (delete people::PendingUserRequest filter .email = <str>$0) { ** }
		`, &deleted, email,
	)
	return
}

// SendConfirmationEmail sends a confirmation email to the user with a verification token.
// It generates a confirmation token, and sends an email with the confirmation link.
// The confirmation token is included as a query parameter in the URL.
func (p *PendingUserRequest) SendConfirmationEmail(db *edgedb.Client, target url.URL) error {
	emailToken := tokens.NewEmailVerificationToken(p.Email)

	if err := emailToken.Save(db); err != nil {
		return err
	}

	params := target.Query()
	params.Set("token", string(emailToken.Token))
	target.RawQuery = params.Encode()

	templateData := email_templates.EmailVerificationData{
		Name: p.FirstName,
		URL:  target,
	}

	return (&email.EmailData{
		To:       emailToken.Email,
		From:     settings.Email().FromHeader(),
		Subject:  templateData.Subject(),
		Template: email_templates.EmailVerification(templateData),
	}).Send(settings.Email().FromHeader())
}

// VerifyEmail attempts to match a token to an EmailVerification entry
// in the database.
// If successful, the token is consumed and the associated account request
// is marked as verified.
func VerifyEmail(edb *edgedb.Client, token tokens.Token) (ok bool, err error) {
	db_token, err := tokens.RetrieveEmailToken(edb, token)
	if err != nil {
		// Token not found is just an invalid token
		if db.IsNoData(err) {
			return false, nil
		}
		return false, err
	}

	if !db_token.IsValid() {
		return false, nil
	}

	// Consume token and set email verified
	txErr := edb.Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
		pending_user, err := GetPendingUserRequest(edb, db_token.Email)
		if err != nil {
			return err
		}
		if err := pending_user.SetEmailVerified(edb, true); err != nil {
			return err
		}
		if err := db_token.Consume(edb); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return false, txErr
	}

	// Email successfully verified
	return true, nil
}

type SuperAdminInput struct {
	UserInput      `edgedb:"$inline" json:",inline"`
	PersonIdentity `edgedb:"$inline" json:",inline"`
	Alias          models.OptionalInput[string] `json:"alias,omitempty" fake:"-"`
	Organisation   OrganisationInput            `edgedb:"organisation" json:"organisation"`
}

func (i SuperAdminInput) Save(e edgedb.Executor) (created User, err error) {
	data, _ := json.Marshal(i)
	if !i.Alias.IsSet {
		i.Alias.Value = i.GenerateAlias()
	}
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with module people,
			data := <json>$0,
			user := (insert User {
				login := <str>data['login'],
				email := <str>data['email'],
				password := <str>data['password'],
				role := UserRole.Admin,
				identity := (insert Person {
					first_name := <str>data['first_name'],
					last_name := <str>data['last_name'],
					contact := <str>data['email'],
					alias := <str>json_get(data, 'alias') ?? {},
					organisations := (insert people::Organisation {
						name := <str>data['organisation']['name'],
						code := <str>data['organisation']['code'],
						description := <str>json_get(data['organisation'], 'description'),
						kind := <OrgKind>data['organisation']['kind']
					})
				})
			}),
			select user { ** }
		`, &created, data)
	return
}
