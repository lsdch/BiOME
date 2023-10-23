package users

import (
	"context"
	"darco/proto/models"
	_ "embed"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string // @name UserRole

const (
	Guest         UserRole = "Guest"
	Contributor   UserRole = "Contributor"
	ProjectMember UserRole = "ProjectMember"
	Admin         UserRole = "Admin"
)

func (m *UserRole) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(*m), nil
}

func (m *UserRole) UnmarshalEdgeDBStr(data []byte) error {
	*m = UserRole(string(data))
	return nil
}

type InnerUserInput struct {
	Login       string      `edgedb:"login" json:"login" validate:"alphanum"`
	Email       string      `edgedb:"email" json:"email" validate:"email" format:"email"`
	EmailPublic bool        `edbedb:"email_public" json:"email_public"`
	Person      PersonInput `edgedb:"identity" json:"identity"`
}

type PasswordInput struct {
	Password   string `json:"password" validate:"required,gte=8"`
	ConfirmPwd string `json:"password_confirmation" validate:"eqfield=Password,required"`
} //@name PasswordInput

func (pwdInput *PasswordInput) Hash() {}

type UserInput struct {
	InnerUserInput `json:",inline"`
	PasswordInput  `json:",inline"`
} // @name UserInput

func (input *UserInput) ProcessPassword() (*UserInsert, error) {
	hashed_password, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	return &UserInsert{
		InnerUserInput: input.InnerUserInput,
		PasswordHash:   hashed_password,
	}, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

type UserInsert struct {
	InnerUserInput `edgedb:"$inline" json:",inline"`
	PasswordHash   string `edgedb:"password" json:"password"`
}

type UserPartial struct {
	Role     UserRole `edgedb:"role" json:"role"`
	Verified bool     `edgedb:"verified" json:"verified"`
	Person   Person   `edgedb:"identity" json:"identity"`
} // @name UserPartial

type User struct {
	UserPartial `edgedb:"$inline" json:",inline"`
	ID          uuid.UUID `edgedb:"id" json:"-"`
	Email       string    `edgedb:"email" json:"email"`
	Login       string    `edgedb:"login" json:"login"`
	Password    string    `edgedb:"password" json:"-"`
} //@name User

func (user *User) Partial() *UserPartial {
	return &user.UserPartial
}

func (user *User) SetActive(active bool) error {
	query := `update people::User
	filter .email = <str>$1
	set {
		verified := <bool>$0
	} `
	return models.DB.Execute(context.Background(), query, active, user.Email)
}

func (user *User) MarshalJSON() ([]byte, error) {
	u := User(*user)
	u.Password = "**********"
	return json.Marshal(u)
}

type UserCredentials struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Remember   bool   `json:"remember" validate:"required"`
} // @name UserCredentials

func (creds *UserCredentials) Authenticate() (*User, error) {
	authFailedError := errors.New("invalid user identifier or password")
	user, err := Find(creds.Identifier)
	if err != nil {
		return nil, authFailedError
	}

	if !user.Verified {
		return nil, errors.New("account is not verified")
	}

	if err := VerifyPassword(user.Password, creds.Password); err != nil {
		return nil, authFailedError
	}
	return user, nil
}

// Find a user by login or email
func Find(identifier string) (*User, error) {
	var user User

	query := `select people::User {
		login, email, verified, role, password, identity: { * }
	} filter .email = <str>$0 or .login = <str>$0 limit 1`
	if err := models.DB.QuerySingle(context.Background(), query, &user, identifier); err != nil {
		return nil, err
	}
	return &user, nil
}
