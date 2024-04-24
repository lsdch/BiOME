package accounts

import (
	"context"
	"darco/proto/db"
	users "darco/proto/models/people"
	"darco/proto/resolvers"
	"darco/proto/services/tokens"
	"fmt"
	"net/http"
	"slices"

	"github.com/sirupsen/logrus"
)

type AuthTokenResponse struct {
	Token string `json:"token" doc:"JSON Web Token" example:"xxxxx.yyyyy.zzzzz"`
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
				AuthTokenResponse: AuthTokenResponse{Token: input.AuthToken},
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
	User              *users.User
}

type LoginOutput struct {
	SessionCookie http.Cookie `header:"Set-Cookie" doc:"Session cookie storing JWT"`
	Body          AuthenticationResponse
}

func createSession(user *users.User, domain string, messages ...string) (*LoginOutput, error) {
	logrus.Infof("Starting user session for: %s [%v]", user.Person.FullName, user.ID)
	token, err := user.GenerateJWT()
	if err != nil {
		respError := fmt.Errorf("Failed to generate session token: %w", err)
		logrus.Error(respError)
		return nil, respError
	}

	cookie := user.JWTCookie(token, domain)
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
			AuthTokenResponse: AuthTokenResponse{Token: token},
			User:              user,
		},
	}, nil
}

func Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user, authError := input.Body.Authenticate(db.Client())
	if authError != nil {
		return nil, authError
	}
	return createSession(user, input.Host)
}

type LogoutOutput struct {
	SetCookie []*http.Cookie `header:"Set-Cookie"`
}

func Logout(ctx context.Context, input *struct{}) (*LogoutOutput, error) {
	return &LogoutOutput{
		SetCookie: []*http.Cookie{
			{Name: tokens.AUTH_TOKEN_COOKIE, MaxAge: -1},
			{Name: tokens.REFRESH_TOKEN_COOKIE, MaxAge: -1},
		},
	}, nil
}
