package accounts

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	users "github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/tokens"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/services/auth_tokens"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirupsen/logrus"
)

type AuthTokenResponse struct {
	AuthToken        string    `json:"auth_token" doc:"JSON Web Token" example:"xxxxx.yyyyy.zzzzz"`
	RefreshToken     string    `json:"refresh_token" doc:"Session refresh token"`
	AuthTokenExpires time.Time `json:"auth_token_expiration" doc:"Time at which auth token expires"`
}

type CurrentUserInput struct {
	resolvers.AuthResolver
}

type CurrentUserResponse struct {
	User              *users.User `json:"user"`
	AuthTokenResponse `json:",inline"`
}

type CurrentUserOutput struct {
	Status int
	Body   CurrentUserResponse `json:"omitempty"`
}

func CurrentUser(ctx context.Context, input *CurrentUserInput) (*CurrentUserOutput, error) {
	if input.User != nil {
		return &CurrentUserOutput{
			Status: http.StatusOK,
			Body: CurrentUserResponse{
				User:              input.User,
				AuthTokenResponse: AuthTokenResponse{AuthToken: input.AuthToken},
			},
		}, nil
	} else {
		return &CurrentUserOutput{
			Status: http.StatusNoContent,
		}, nil
	}
}

type LoginInput struct {
	resolvers.HostResolver
	Body users.UserCredentials `required:"true"`
}

type AuthenticationResponse struct {
	Messages          []string `json:"messages"`
	AuthTokenResponse `json:",inline"`
	User              *users.User `json:"user"`
}

type LoginOutput struct {
	SessionCookie http.Cookie `header:"Set-Cookie" doc:"Session cookie storing JWT"`
	Body          AuthenticationResponse
}

type SessionParameters struct {
	Domain       string
	RefreshToken tokens.SessionRefreshToken
}

func createSession(db geltypes.Executor, params SessionParameters, messages ...string) (*LoginOutput, error) {
	user, err := users.FindID(db, params.RefreshToken.UserID)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Starting user session for: %s [%v]", user.Person.FullName, user.Role)
	token, err := user.GenerateJWT()
	if err != nil {
		respError := fmt.Errorf("Failed to generate session token: %v", err)
		logrus.Error(respError)
		return nil, respError
	}

	cookie := user.JWTCookie(string(token.Token), params.Domain)
	logrus.Debugf("Set JWT session cookie: %+v", cookie)

	return &LoginOutput{
		SessionCookie: cookie,
		Body: AuthenticationResponse{
			Messages: slices.Concat(
				[]string{fmt.Sprintf(
					"User authenticated as %s %s",
					user.Person.FirstName, user.Person.LastName,
				)},
				messages,
			),
			AuthTokenResponse: AuthTokenResponse{
				AuthToken:        string(token.Token),
				RefreshToken:     string(params.RefreshToken.Token),
				AuthTokenExpires: token.Expires,
			},
			User: &user,
		},
	}, nil
}

func Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user, authError := input.Body.Authenticate(db.Client())
	if db.IsNoData(authError) {
		return nil, huma.Error401Unauthorized("Invalid credentials")
	}
	if authError != nil {
		return nil, authError
	}
	refreshToken, err := tokens.CreateSessionRefreshToken(db.Client(), user.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate refresh token: %v", err)
	}
	return createSession(db.Client(),
		SessionParameters{
			Domain:       input.HostResolver.Host,
			RefreshToken: refreshToken,
		})
}

type LogoutInput struct {
	Body struct {
		RefreshToken models.OptionalInput[string] `json:"refresh_token,omitempty"`
	}
}

type LogoutOutput struct {
	SetCookie []*http.Cookie `header:"Set-Cookie"`
}

func Logout(ctx context.Context, input *LogoutInput) (*LogoutOutput, error) {
	token, ok := input.Body.RefreshToken.Get()
	if ok {
		_ = tokens.DropSessionToken(db.Client(), tokens.Token(token))
	}
	return &LogoutOutput{
		SetCookie: []*http.Cookie{
			{Name: auth_tokens.AUTH_TOKEN_COOKIE, MaxAge: -1, Value: "", Path: "/", HttpOnly: true, Secure: true},
			{Name: auth_tokens.REFRESH_TOKEN_COOKIE, MaxAge: -1, Value: "", Path: "/", HttpOnly: true, Secure: true},
		},
	}, nil
}

type RefreshTokenBody struct {
	Token string `json:"refresh_token"`
}

type RefreshSessionInput struct {
	resolvers.HostResolver
	Body RefreshTokenBody
}

func RefreshSession(ctx context.Context, input *RefreshSessionInput) (*LoginOutput, error) {
	token, err := tokens.RetrieveSessionRefreshToken(
		db.Client(),
		tokens.Token(input.Body.Token),
	)
	if db.IsNoData(err) || !token.IsValid() {
		return nil, huma.Error401Unauthorized("Invalid token")
	}
	if err != nil {
		return nil, err
	}
	newToken, err := token.Rotate(db.Client())
	if err != nil {
		return nil, err
	}

	return createSession(db.Client(),
		SessionParameters{
			Domain:       input.Host,
			RefreshToken: newToken,
		},
	)
}
