package accounts

import (
	"context"
	"errors"
	"fmt"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/people"
	users "github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/tokens"
	_ "github.com/lsdch/biome/models/validations"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirupsen/logrus"
)

type RegisterInput struct {
	resolvers.HostResolver
	Body struct {
		Data             people.PendingUserRequestInput `json:"data"`
		VerificationPath string                         `json:"verification_path"`
	}
}

func Register(confirmEmailPath string) router.Endpoint[RegisterInput, controllers.Message] {
	return func(ctx context.Context, input *RegisterInput) (*controllers.Message, error) {
		logrus.Infof("Attempting to create account for %s %s (%s)",
			input.Body.Data.FirstName,
			input.Body.Data.LastName,
			input.Body.Data.Email,
		)
		err := db.Client().Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {
			pending, err := input.Body.Data.Register(db.Client())
			if err != nil {
				return fmt.Errorf("Failed to create account request: %v", err)
			}

			target := input.OriginPath(input.Body.VerificationPath)
			if err := pending.SendConfirmationEmail(db.Client(), target); err != nil {
				return fmt.Errorf("Failed to send verification email: %v", err)
			}
			return nil
		})

		return &controllers.Message{
			Body: "Account request created and email with verification token was sent",
		}, err

	}
}

type ConfirmEmailInput struct {
	resolvers.HostResolver
	Token string `query:"token"`
}

func ConfirmEmail(ctx context.Context, input *ConfirmEmailInput) (*struct{ Message string }, error) {

	ok, err := people.VerifyEmail(db.Client(), tokens.Token(input.Token))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, huma.Error422UnprocessableEntity("Token is invalid or expired")
	}
	return &struct{ Message string }{"Email successfully verified"}, nil
}

type ResendEmailVerificationInput struct {
	resolvers.HostResolver
	Body struct {
		Email                string               `json:"email" format:"email"`
		EmailVerificationURL TokenVerificationURL `json:"verification_url"`
	}
}

func ResendEmailVerification(confirmEmailPath string) router.Endpoint[ResendEmailVerificationInput, struct{}] {
	return func(ctx context.Context, input *ResendEmailVerificationInput) (*struct{}, error) {
		target := input.GenerateURL(confirmEmailPath)
		if input.Body.EmailVerificationURL.Host != "" {
			target = input.Body.EmailVerificationURL.URL()
		}
		pending, err := users.GetPendingUserRequest(db.Client(), input.Body.Email)
		if err != nil {
			return nil, nil
		}
		if err := pending.SendConfirmationEmail(db.Client(), target); err != nil {
			return nil, huma.Error500InternalServerError("Failed to send verification email", err)
		}
		return nil, nil
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
	refreshToken, err := tokens.CreateSessionRefreshToken(db.Client(), user.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate refresh token: %v", err)
	}
	return createSession(db.Client(),
		SessionParameters{
			Domain:       input.Host,
			RefreshToken: refreshToken,
		},
		fmt.Sprintf("Account created with role %s", user.Role),
	)
}
