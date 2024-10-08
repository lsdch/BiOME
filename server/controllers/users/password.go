package accounts

import (
	"context"
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/models/tokens"
	"darco/proto/resolvers"

	"github.com/danielgtaylor/huma/v2"
)

type ValidatePasswordTokenInput struct {
	Token tokens.Token `query:"token"`
}

func ValidatePasswordToken(ctx context.Context, input *ValidatePasswordTokenInput) (*struct{}, error) {
	token, err := tokens.RetrievePwdResetToken(
		db.Client(),
		tokens.Token(input.Token),
	)
	if db.IsNoData(err) || !token.IsValid() {
		return nil, huma.Error422UnprocessableEntity("Invalid token")
	}
	if err != nil {
		return nil, err
	}

	// Respond with HTTP 204: token is valid
	return nil, nil
}

type UpdatePasswordInput struct {
	resolvers.AuthRequired
	Body people.UpdatePasswordInput `required:"true"`
}

func (i *UpdatePasswordInput) Resolve(ctx huma.Context) []error {
	errs := i.AuthRequired.Resolve(ctx)
	if errs != nil {
		return errs
	}
	// Verify current password
	if !i.User.PasswordMatch(i.DB(), i.Body.Password) {
		errs = append(errs, &huma.ErrorDetail{
			Message:  "Invalid password",
			Location: "password",
			Value:    i.Body.Password,
		})
	}
	// Check password strength
	if !i.User.ValidatePasswordStrength(
		i.DB(), i.Body.NewPassword.Password,
	) {
		errs = append(errs, &huma.ErrorDetail{
			Message:  "Password is too weak",
			Location: "new_password.password",
			Value:    i.Body.NewPassword.Password,
		})
	}
	// Check password confirmation
	if !i.Body.NewPassword.ValidateEqual() {
		errs = append(errs, &huma.ErrorDetail{
			Message:  "Passwords do not match",
			Location: "new_password.password_confirmation",
			Value:    i.Body.NewPassword.ConfirmPwd,
		})
	}

	return errs
}

func UpdatePassword(ctx context.Context, input *UpdatePasswordInput) (*struct{}, error) {
	err := input.User.SetPassword(input.DB(), input.Body.NewPassword.Password)
	return nil, err
}

type PasswordResetInput struct {
	Token tokens.Token         `query:"token"`
	Body  people.PasswordInput `required:"true"`
	*people.User
}

func (i *PasswordResetInput) Resolve(ctx huma.Context) []error {

	resetToken, err := tokens.RetrievePwdResetToken(db.Client(), i.Token)
	if db.IsNoData(err) || !resetToken.IsValid() {
		return []error{&huma.ErrorDetail{Message: "Invalid token"}}
	}
	if err != nil {
		return []error{err}
	}
	user, err := people.FindID(db.Client(), resetToken.UserID)
	if err != nil {
		return []error{err}
	}
	i.User = &user
	return nil
}

func PasswordReset(ctx context.Context, input *PasswordResetInput) (*struct{}, error) {
	return nil, input.User.SetPassword(db.Client(), input.Body.Password)
}

type RequestPasswordResetInput struct {
	resolvers.HostResolver
	Body struct {
		Email   string               `json:"email" format:"email"`
		Handler TokenVerificationURL `json:"handler,omitempty" doc:"A URL where a form to set the new password is available"`
	} `nameHint:"PasswordResetRequest"`
}

type RequestPasswordResetHandler func(context.Context, *RequestPasswordResetInput) (*struct{}, error)

func RequestPasswordReset(defaultHandlerPath string) RequestPasswordResetHandler {
	return func(ctx context.Context, input *RequestPasswordResetInput) (*struct{}, error) {
		user, err := people.Find(db.Client(), input.Body.Email)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("Unknown email address", err)
		}

		var targetURL = input.GenerateURL(defaultHandlerPath)
		if input.Body.Handler.Host == "" {
			targetURL = input.Body.Handler.URL()
		}

		if err = user.RequestPasswordReset(db.Client(), targetURL); err != nil {
			return nil, huma.Error500InternalServerError(
				"Failed to issue password reset token",
				err,
			)
		}
		return nil, nil
	}
}
