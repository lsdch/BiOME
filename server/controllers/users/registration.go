package accounts

import (
	"context"
	"darco/proto/controllers"
	"darco/proto/db"
	"darco/proto/models/people"
	users "darco/proto/models/people"
	"darco/proto/models/tokens"
	_ "darco/proto/models/validations"
	"darco/proto/resolvers"
	"darco/proto/router"
	"errors"
	"fmt"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirupsen/logrus"
)

type EmailVerificationURL struct {
	Handler url.URL `json:"handler,omitempty" doc:"A URL used to generate the verification link, which can be set by the web client. Verification token will be added as a URL query parameter."`
}

type RegisterInput struct {
	resolvers.HostResolver
	Body struct {
		people.PendingUserRequestInput `json:",inline"`
		EmailVerificationURL           `json:",inline"`
	}
}

func Register(confirmEmailPath string) router.Endpoint[RegisterInput, controllers.Message] {
	return func(ctx context.Context, input *RegisterInput) (*controllers.Message, error) {
		logrus.Infof("Attempting to create account for %s %s (%s)",
			input.Body.Person.FirstName,
			input.Body.Person.LastName,
			input.Body.Email,
		)

		pending, err := input.Body.Register(db.Client())
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to register new account", err)
		}

		target := input.GenerateURL(confirmEmailPath)
		if input.Body.Handler.Host != "" {
			target = input.Body.Handler
		}
		if err := pending.SendConfirmationEmail(db.Client(), target); err != nil {
			return nil, huma.Error500InternalServerError("Failed to send verification email", err)
		}
		return &controllers.Message{
			Body: "Account request created and email with verification token was sent",
		}, nil

	}
}

type ConfirmEmailInput struct {
	resolvers.HostResolver
	Token string `query:"token"`
	*users.PendingUserRequest
}

func (i *ConfirmEmailInput) Resolve(ctx huma.Context) []error {
	if errs := i.HostResolver.Resolve(ctx); errs != nil {
		return errs
	}

	token, err := tokens.RetrieveEmailToken(db.Client(), tokens.Token(i.Token))
	if db.IsNoData(err) || !token.IsValid() {
		return []error{&huma.ErrorDetail{Message: "Invalid token"}}
	}
	if err != nil {
		return []error{err}
	}

	accountRequest, err := people.GetPendingUserRequest(db.Client(), token.Email)
	if db.IsNoData(err) {
		return []error{fmt.Errorf("No pending account request associated to this email found")}
	}

	i.PendingUserRequest = accountRequest

	return nil
}

func ConfirmEmail(ctx context.Context, input *ConfirmEmailInput) (*struct{ Message string }, error) {

	if err := input.PendingUserRequest.SetEmailVerified(db.Client(), true); err != nil {
		return nil, huma.Error500InternalServerError("Email confirmation failed", err)
	}
	return &struct{ Message string }{"Email successfully verified"}, nil
}

type ResendEmailConfirmationInput struct {
	resolvers.HostResolver
	Body struct {
		Email                string `json:"email" format:"email"`
		EmailVerificationURL `json:",inline"`
	}
	*users.PendingUserRequest
}

func (i *ResendEmailConfirmationInput) Resolve(ctx huma.Context) []error {
	pending, err := users.GetPendingUserRequest(db.Client(), i.Body.Email)
	if err != nil {
		return []error{&huma.ErrorDetail{
			Message:  "Unknown e-mail address",
			Location: "email",
			Value:    i.Body.Email,
		}}
	}
	if pending.EmailVerified {
		return []error{&huma.ErrorDetail{
			Message:  "E-mail was already verified",
			Location: "email",
			Value:    i.Body.Email,
		}}
	}
	i.PendingUserRequest = pending
	return nil
}

func ResendEmailConfirmation(confirmEmailPath string) router.Endpoint[ResendEmailConfirmationInput, controllers.Message] {
	return func(ctx context.Context, input *ResendEmailConfirmationInput) (*controllers.Message, error) {
		target := input.GenerateURL(confirmEmailPath)
		if input.Body.Handler.Host != "" {
			target = input.Body.Handler
		}
		if err := input.PendingUserRequest.SendConfirmationEmail(db.Client(), target); err != nil {
			return nil, huma.Error500InternalServerError("Failed to send verification email", err)
		}
		return &controllers.Message{Body: "Verification email was sent"}, nil
	}
}

type ClaimInvitationInput struct {
	resolvers.HostResolver
	Token tokens.Token `path:"token"`
	Body  people.UserInput
}

func ClaimInvitation(ctx context.Context, input *ClaimInvitationInput) (*LoginOutput, error) {
	user, err := input.Body.RegisterWithToken(db.Client(), input.Token)
	switch {
	case errors.Is(err, people.InvalidTokenError):
		return nil, huma.Error422UnprocessableEntity("Invalid invitation token")
	case err != nil:
		return nil, huma.Error500InternalServerError("Registration failed due to an internal server error", err)
	}
	return createSession(user, input.Host,
		fmt.Sprintf("Account created with role %s", user.Role),
	)
}
