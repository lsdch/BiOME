package accounts

import (
	"darco/proto/controllers"
	"darco/proto/models/people"
	"darco/proto/resolvers"
	"darco/proto/router"
	"fmt"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

// Path to the API endpoint where invitation token can be consumed to register an account.
// See [person.InvitePerson].
var InvitationClaimPath = ""

func RegisterRoutes(r router.Router) {
	accountAPI := r.RouteGroup("/account").
		WithTags([]string{"Account"})

	registry := r.API.OpenAPI().Components.Schemas

	router.Register(accountAPI, "CurrentUser",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "Current user",
			Description: "Get infos of currently authenticated user account",
			Responses: map[string]*huma.Response{
				"200": {
					Description: "The currently authenticated user",
					Content: map[string]*huma.MediaType{
						"application/json": {
							Schema: registry.Schema(reflect.TypeFor[CurrentUserResponse](), true, ""),
						},
					},
				},
				"204": {Description: "No active user session", Content: nil},
			},
		}, CurrentUser)

	router.Register(accountAPI, "Login",
		huma.Operation{
			Path:        "/login",
			Method:      http.MethodPost,
			Summary:     "Login",
			Description: "Authenticate using user credentials",
			Errors:      []int{http.StatusUnprocessableEntity},
		}, Login)

	router.Register(accountAPI, "Logout",
		huma.Operation{
			Path:        "/logout",
			Method:      http.MethodPost,
			Summary:     "Logout",
			Description: "Logout from current user session by revoking session cookies",
		}, Logout)

	router.Register(accountAPI, "RefreshSession",
		huma.Operation{
			Path:        "/refresh",
			Method:      http.MethodPost,
			Summary:     "Refresh auth token",
			Description: "Refresh session using refresh token",
		}, RefreshSession)

	router.Register(accountAPI, "UpdatePassword",
		huma.Operation{
			Path:        "/password",
			Method:      http.MethodPost,
			Summary:     "Update password",
			Description: "Updates password of currently authenticated user",
			Errors:      []int{http.StatusUnprocessableEntity, http.StatusUnauthorized},
		}, UpdatePassword)

	pwdResetPath := router.Register(accountAPI, "ResetPassword",
		huma.Operation{
			Path:        "/password-reset/",
			Method:      http.MethodPost,
			Summary:     "Reset password",
			Description: "Set a new password using a previously issued reset token",
		},
		PasswordReset)

	router.Register(accountAPI, "ValidatePasswordToken",
		huma.Operation{
			Path:        "password-reset/",
			Method:      http.MethodGet,
			Summary:     "Validate password token",
			Description: "Verifies that the password token is valid and can be used to reset a password",
			Errors:      []int{http.StatusUnprocessableEntity},
		}, ValidatePasswordToken)

	router.Register(accountAPI, "RequestPasswordReset",
		huma.Operation{
			Path:        "/forgotten-password",
			Method:      http.MethodPost,
			Summary:     "Request password reset",
			Description: fmt.Sprintf("Requests sending a link containing a password reset token to your account email address. The link target can be provided by the client in the request body, or defaults to the API endpoint: `%s`. In this case, setting the new password is expected to be done programatically, e.g. through a curl request.", pwdResetPath),
			Errors:      []int{http.StatusUnprocessableEntity, http.StatusInternalServerError},
		}, RequestPasswordReset(pwdResetPath))

	confirmEmailPath := router.Register(accountAPI, "ConfirmEmail",
		huma.Operation{
			Path:        "/email-confirmation",
			Method:      http.MethodGet,
			Summary:     "Confirm e-mail",
			Description: "Confirms the validity of an e-mail address associated to an account, using a token issued at the end of user registration.",
			Errors:      []int{http.StatusUnprocessableEntity, http.StatusInternalServerError},
		}, ConfirmEmail)

	router.Register(accountAPI, "ResendEmailConfirmation",
		huma.Operation{
			Path:        "/email-confirmation/resend",
			Method:      http.MethodPost,
			Summary:     "Resend e-mail verification link",
			Description: "Sends again a verification link for the provided e-mail address, if it matches a currently not verified user account.",
			Errors:      []int{http.StatusUnprocessableEntity, http.StatusInternalServerError},
		}, ResendEmailConfirmation(confirmEmailPath))

	router.Register(accountAPI, "ListPendingUserRequests",
		huma.Operation{
			Path:        "/pending",
			Method:      http.MethodGet,
			Summary:     "List pending user requests",
			Description: "Lists all account requests pending validation from an administrator",
		}, controllers.ListHandler[*struct {
			resolvers.AccessRestricted[resolvers.Admin]
		}](people.ListPendingUserRequests))

	router.Register(accountAPI, "GetPendingUserRequest",
		huma.Operation{
			Path:        "/pending/{email}",
			Method:      http.MethodGet,
			Summary:     "Get pending user request",
			Description: "Get account request pending validation using the associated email",
		}, controllers.GetHandler[*struct {
			resolvers.AccessRestricted[resolvers.Admin]
			controllers.StrIdentifier `path:"email" format:"email"`
		}](people.GetPendingUserRequest))

	router.Register(accountAPI, "DeletePendingUserRequest",
		huma.Operation{
			Path:        "/pending/{email}",
			Method:      http.MethodDelete,
			Summary:     "Delete pending user request",
			Description: "Delete account request pending validation using the associated email",
		}, controllers.DeleteHandler[*struct {
			controllers.StrIdentifier `path:"email" format:"email"`
			resolvers.AccessRestricted[resolvers.Admin]
		}](people.DeletePendingUserRequest))

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
		}, Register(confirmEmailPath))

	InvitationClaimPath = router.Register(accountAPI, "ClaimInvitation",
		huma.Operation{
			Path:        "/register/{token}",
			Method:      http.MethodPost,
			Summary:     "Claim invitation",
			Description: "Register an account with pre-assigned role and identity, using an invitation token",
		}, ClaimInvitation)
}
