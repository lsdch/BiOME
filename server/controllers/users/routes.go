package accounts

import (
	"darco/proto/router"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	accountAPI := r.RouteGroup("/account").
		WithTags([]string{"Account"})

	router.Register(accountAPI, "Login",
		huma.Operation{
			Path:        "/login",
			Method:      http.MethodPost,
			Summary:     "Login",
			Description: "Authenticate using user credentials",
			Errors:      []int{400},
		}, Login)

	router.Register(accountAPI, "Logout",
		huma.Operation{
			Path:        "/logout",
			Method:      http.MethodPost,
			Summary:     "Logout",
			Description: "Logout from current user session by revoking session cookies",
		}, Logout)

	router.Register(accountAPI, "CurrentUser",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "Current user",
			Description: "Get infos of currently authenticated user account",
			Errors:      []int{404},
		}, CurrentUser)

	router.Register(accountAPI, "UpdatePassword",
		huma.Operation{
			Path:        "/password",
			Method:      http.MethodPost,
			Summary:     "Update password",
			Description: "Updates password of currently authenticated user",
			Errors:      []int{400, 401},
		}, UpdatePassword)

	pwdResetPath := router.Register(accountAPI, "ResetPassword",
		huma.Operation{
			Path:        "/password-reset/{token}",
			Method:      http.MethodPost,
			Summary:     "Reset password",
			Description: "Set a new password using a previously issued reset token",
			Errors:      []int{400},
		},
		PasswordReset)

	router.Register(accountAPI, "ValidatePasswordToken",
		huma.Operation{
			Path:        "password-reset/{token}",
			Method:      http.MethodGet,
			Summary:     "Validate password token",
			Description: "Verifies that the password token is valid and can be used to reset a password",
			Errors:      []int{400},
		}, ValidatePasswordToken)

	router.Register(accountAPI, "RequestPasswordReset",
		huma.Operation{
			Path:        "/forgotten-password",
			Method:      http.MethodPost,
			Summary:     "Request password reset",
			Description: fmt.Sprintf("Requests sending a link containing a password reset token to your account email address. The link target can be provided by the client in the request body, or defaults to the API endpoint: `%s`. In this case, setting the new password is expected to be done programatically, e.g. through a curl request.", pwdResetPath),
			Errors:      []int{400, 500},
		}, RequestPasswordReset(pwdResetPath))

	confirmEmailPath := router.Register(accountAPI, "ConfirmEmail",
		huma.Operation{
			Path:        "/email-confirmation",
			Method:      http.MethodGet,
			Summary:     "Confirm e-mail",
			Description: "Confirms the validity of an e-mail address associated to an account, using a token issued at the end of user registration.",
			Errors:      []int{400, 500},
		}, ConfirmEmail)

	router.Register(accountAPI, "ResendEmailConfirmation",
		huma.Operation{
			Path:        "/email-confirmation/resend",
			Method:      http.MethodPost,
			Summary:     "Resend e-mail verification link",
			Description: "Sends again a verification link for the provided e-mail address, if it matches a currently not verified user account.",
			Errors:      []int{400, 500},
		}, ResendEmailConfirmation(confirmEmailPath))

	router.Register(accountAPI, "Register",
		huma.Operation{
			Path:    "/register",
			Method:  http.MethodPost,
			Summary: "Register new account",
			Description: fmt.Sprintf(
				"Register a new account that is initially pending, and needs to be activated by an administrator. An email is sent to the registered e-mail address with a verification link. The target URL can be set by the client, otherwise it defaults to the API endpoint: `%s`",
				confirmEmailPath,
			),
			DefaultStatus: http.StatusCreated,
			Errors:        []int{400, 500},
		}, Register(confirmEmailPath))
}
