package accounts

import (
	"context"
	"darco/proto/db"
	"darco/proto/models/people"
	users "darco/proto/models/people"
	_ "darco/proto/models/validations"
	"darco/proto/resolvers"
	"darco/proto/router"
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

func sendConfirmationEmail(user *users.User, target url.URL) (*struct{}, error) {
	if err := user.SendConfirmationEmail(db.Client(), target); err != nil {
		msg := fmt.Sprintf(
			"Failed to send account confirmation email to '%s'.",
			user.Email,
		)
		logrus.Errorf("%s Error: %v", msg, err)
		return nil, huma.Error500InternalServerError(msg, err)
	}

	return nil, nil
}

func Register(confirmEmailPath string) router.Endpoint[RegisterInput, struct{}] {
	return func(ctx context.Context, input *RegisterInput) (*struct{}, error) {
		logrus.Infof("Attempting to create account for %s %s (%s)",
			input.Body.Person.FirstName,
			input.Body.Person.LastName,
			input.Body.User.Email,
		)

		pending, err := input.Body.Register(db.Client())
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to register new account", err)
		}

		target := input.GenerateURL(confirmEmailPath)
		if input.Body.Handler.Host != "" {
			target = input.Body.Handler
		}
		return sendConfirmationEmail(&pending.User, target)
	}
}

type ConfirmEmailInput struct {
	resolvers.HostResolver
	Token string `query:"token"`
	*users.User
}

func (i *ConfirmEmailInput) Resolve(ctx huma.Context) []error {
	if errs := i.HostResolver.Resolve(ctx); errs != nil {
		return errs
	}
	user, tokenValid := users.ValidateAccountToken(
		db.Client(),
		users.Token(i.Token),
		users.EmailConfirmationToken,
	)
	if !tokenValid {
		return []error{&huma.ErrorDetail{
			Message: "Invalid token",
		}}
	}

	if user.EmailConfirmed {
		return []error{&huma.ErrorDetail{
			Message: "Account is already verified",
		}}
	}

	i.User = user
	return nil
}

func ConfirmEmail(ctx context.Context, input *ConfirmEmailInput) (*LoginOutput, error) {
	if err := input.User.SetEmailConfirmed(db.Client(), true); err != nil {
		return nil, huma.Error500InternalServerError("Email confirmation failed", err)
	}
	return createSession(
		input.User,
		input.Host,
		"Email confirmation successful",
	)
}

type ResendEmailConfirmationInput struct {
	resolvers.HostResolver
	Body struct {
		Email                string `json:"email" format:"email"`
		EmailVerificationURL `json:",inline"`
	}
	*users.User
}

func (i *ResendEmailConfirmationInput) Resolve(ctx huma.Context) []error {
	user, err := users.Find(db.Client(), i.Body.Email)
	if err != nil {
		return []error{&huma.ErrorDetail{
			Message:  "Unknown e-mail address",
			Location: "email",
			Value:    i.Body.Email,
		}}
	}
	if user.EmailConfirmed {
		return []error{&huma.ErrorDetail{
			Message:  "E-mail was already verified",
			Location: "email",
			Value:    i.Body.Email,
		}}
	}
	i.User = user
	return nil
}

func ResendEmailConfirmation(confirmEmailPath string) router.Endpoint[ResendEmailConfirmationInput, struct{}] {
	return func(ctx context.Context, input *ResendEmailConfirmationInput) (*struct{}, error) {
		target := input.GenerateURL(confirmEmailPath)
		if input.Body.Handler.Host != "" {
			target = input.Body.Handler
		}
		return sendConfirmationEmail(input.User, target)
	}
}
