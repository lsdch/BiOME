package accounts

import (
	"fmt"
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

// Path to the API endpoint where invitation token can be consumed to register an account.
// See [person.InvitePerson].
var InvitationClaimPath = ""

func init() {
	accountAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/account").
			WithTags([]string{"Account"})
	}

	router.RegisterSpec(
		accountAPI,
		"CurrentUser",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "Current user",
			Description: "Get infos of currently authenticated user account",
		},
		CurrentUser)

	router.RegisterSpec(
		accountAPI,
		"Login",
		huma.Operation{
			Path:        "/login",
			Method:      http.MethodPost,
			Summary:     "Login",
			Description: "Authenticate using user credentials",
			Errors:      []int{http.StatusUnprocessableEntity},
		},
		Login)

	router.RegisterSpec(
		accountAPI,
		"Logout",
		huma.Operation{
			Path:        "/logout",
			Method:      http.MethodPost,
			Summary:     "Logout",
			Description: "Logout from current user session by revoking session cookies",
		},
		Logout)

	router.RegisterSpec(
		accountAPI,
		"RefreshSession",
		huma.Operation{
			Path:        "/refresh",
			Method:      http.MethodPost,
			Summary:     "Refresh auth token",
			Description: "Refresh session using refresh token",
		},
		RefreshSession)

	router.RegisterSpec(
		accountAPI,
		"UpdatePassword",
		huma.Operation{
			Path:        "/password",
			Method:      http.MethodPost,
			Summary:     "Update password",
			Description: "Updates password of currently authenticated user",
			Errors:      []int{http.StatusUnprocessableEntity, http.StatusUnauthorized},
		},
		UpdatePassword)

	pwdResetRoute := router.RegisterSpec(
		accountAPI,
		"ResetPassword",
		huma.Operation{
			Path:        "/password-reset/",
			Method:      http.MethodPost,
			Summary:     "Reset password",
			Description: "Set a new password using a previously issued reset token",
		},
		PasswordReset)

	router.RegisterSpec(
		accountAPI,
		"ValidatePasswordToken",
		huma.Operation{
			Path:        "/password-reset/",
			Method:      http.MethodGet,
			Summary:     "Validate password token",
			Description: "Verifies that the password token is valid and can be used to reset a password",
			Errors:      []int{http.StatusUnprocessableEntity},
		},
		ValidatePasswordToken)

	router.RegisterCustom(func(r *router.Router) {
		router.NewSpec(
			accountAPI,
			"RequestPasswordReset",
			huma.Operation{
				Path:        "/forgotten-password",
				Method:      http.MethodPost,
				Summary:     "Request password reset",
				Description: fmt.Sprintf("Requests sending a link containing a password reset token to your account email address. The link target can be provided by the client in the request body, or defaults to the API endpoint: `%s`. In this case, setting the new password is expected to be done programatically, e.g. through a curl request.", pwdResetRoute.Path(r)),
				Errors:      []int{http.StatusUnprocessableEntity, http.StatusInternalServerError},
			},
			RequestPasswordReset(pwdResetRoute.Path(r))).Register(r)
	})

	confirmEmailRoute := router.RegisterSpec(
		accountAPI,
		"ConfirmEmail",
		huma.Operation{
			Path:        "/email-confirmation",
			Method:      http.MethodGet,
			Summary:     "Confirm e-mail",
			Description: "Confirms the validity of an e-mail address associated to an account, using a token issued at the end of user registration.",
			Errors:      []int{http.StatusUnprocessableEntity, http.StatusInternalServerError},
		},
		ConfirmEmail)

	router.RegisterCustom(func(r *router.Router) {
		router.NewSpec(
			accountAPI,
			"ResendEmailVerification",
			huma.Operation{
				Path:        "/email-confirmation/resend",
				Method:      http.MethodPost,
				Summary:     "Resend e-mail verification link",
				Description: "Sends again a verification link for the provided e-mail address, if it matches a currently not verified user account.",
				Errors:      []int{http.StatusUnprocessableEntity, http.StatusInternalServerError},
			},
			ResendEmailVerification(confirmEmailRoute.Path(r))).Register(r)
	})

	router.RegisterSpec(
		accountAPI,
		"ListPendingUserRequests",
		huma.Operation{
			Path:        "/pending",
			Method:      http.MethodGet,
			Summary:     "List pending user requests",
			Description: "Lists all account requests pending validation from an administrator",
		},
		controllers.ListHandler[*struct {
			resolvers.AccessRestricted[resolvers.Admin]
		}](people.ListPendingUserRequests))

	router.RegisterSpec(
		accountAPI,
		"GetPendingUserRequest",
		huma.Operation{
			Path:        "/pending/{email}",
			Method:      http.MethodGet,
			Summary:     "Get pending user request",
			Description: "Get account request pending validation using the associated email",
		},
		controllers.GetHandler[*struct {
			resolvers.AccessRestricted[resolvers.Admin]
			controllers.StrIdentifier `path:"email" format:"email"`
		}](people.GetPendingUserRequest))

	router.RegisterSpec(
		accountAPI,
		"DeletePendingUserRequest",
		huma.Operation{
			Path:        "/pending/{email}",
			Method:      http.MethodDelete,
			Summary:     "Delete pending user request",
			Description: "Delete account request pending validation using the associated email",
		},
		controllers.DeleteHandler[*struct {
			controllers.EmailInput `path:"email" format:"email"`
			resolvers.AccessRestricted[resolvers.Admin]
		}](people.DeletePendingUserRequest))

	router.RegisterCustom(func(r *router.Router) {
		router.NewSpec(
			accountAPI,
			"Register",
			huma.Operation{
				Path:          "/register",
				Method:        http.MethodPost,
				Summary:       "Register new account",
				Description:   "Register a new account that is initially pending, and needs to be activated by an administrator. An email is sent to the registered e-mail address with a verification link.",
				DefaultStatus: http.StatusCreated,
			},
			Register(confirmEmailRoute.Path(r))).Register(r)
	})
	router.RegisterCustom(func(r *router.Router) {
		// Setting global InvitationClaimPath variable
		InvitationClaimPath = router.NewSpec(
			accountAPI,
			"ClaimInvitation",
			huma.Operation{
				Path:        "/register/{token}",
				Method:      http.MethodPost,
				Summary:     "Claim invitation",
				Description: "Register an account with pre-assigned role and identity, using an invitation token",
			},
			ClaimInvitation).Register(r).Path(r)
	})
}
