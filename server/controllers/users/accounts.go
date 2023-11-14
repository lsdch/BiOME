package accounts

import (
	"context"
	"darco/proto/config"
	"darco/proto/models"
	"darco/proto/models/users"
	_ "darco/proto/models/validations"
	"darco/proto/services/tokens"
	"errors"
	"net/http"
	"net/url"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary Delete a user
// @Description Deletes a user
// @tags Auth
// @Accept json
// @Produce json
// @Success 200 "User was deleted successfully"
// @Failure 401 "Admin privileges required"
// @Failure 404 "User does not exist"
// @Router /users/{uuid} [delete]
func Delete(ctx *gin.Context) {

}

type TokenResponse struct {
	Token string `json:"token" binding:"required" example:"some-generated-jwt"`
} // @name TokenResponse

// @Summary Authenticate user
// @Description Authenticate user with their credentials and set a JWT.
// @id Login
// @tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} TokenResponse "Returns a token and stores it as a session cookie"
// @Failure 400 {object} users.LoginFailedError "Authentication failure"
// @Failure 500
// @Router /login [post]
// @Param data body users.UserCredentials true "User credentials"
func Login(ctx *gin.Context, db *edgedb.Client) {
	var credentials users.UserCredentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			users.LoginFailedError{Reason: users.InvalidCredentials},
		)
		return
	}
	user, authError := credentials.Authenticate(db)
	if authError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, authError)
		return
	}
	startUserSession(ctx, user.ID)
}

// @Summary Logout user
// @Description Log out currently authenticated user
// @id Logout
// @tags Auth
// @Accept json
// @Produce json
// @Success 200 "User logged out"
// @Router /logout [post]
func Logout(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func startUserSession(ctx *gin.Context, userID edgedb.UUID) {
	config := config.Get()
	token, err := tokens.GenerateToken(userID, config.TokenLifetime)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	lifetime := int(config.TokenLifetime.Seconds())
	ctx.SetCookie("token", token, lifetime, "/", "", true, true)
	ctx.JSON(http.StatusOK, TokenResponse{token})
}

// Returns nil if URL parsing failed
func parseRequestOrigin(ctx *gin.Context) *url.URL {
	origin := ctx.Request.Header.Get("Origin")
	originURL, err := url.Parse(origin)
	if err != nil {
		logrus.Errorf("Failed to parse origin URL '%s'. Defaulting to creating token URL pointing to API endpoint.", origin)
		originURL = nil
	}
	return originURL
}

// @Summary Register user
// @Description Register a new user account, that is inactive (until email is verified or admin intervention), and has role 'Guest'
// @id RegisterUser
// @tags Auth
// @Accept json
// @Produce json
// @Success 202 "User created and waiting for email verification"
// @Failure 400 {object} validations.FieldErrors
// @Router /users/register [post]
// @Param data body users.UserInput true "User informations"
func Register(ctx *gin.Context) {
	var newUser users.UserInput
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.Error(err)
		return
	}

	logrus.Infof("Attempting to create account for %s %s (%s)",
		newUser.Person.FirstName,
		newUser.Person.LastName,
		newUser.Email,
	)

	originURL := parseRequestOrigin(ctx)
	if err := newUser.Register(config.Get(), originURL); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusAccepted)
}

type ResendConfirmationError string // @name ResendConfirmationError

const (
	AccountAlreadyVerified   string                  = "AlreadyVerified"
	ResendInvalidCredentials ResendConfirmationError = "InvalidCredentials"
	ResendAlreadyVerified    ResendConfirmationError = ResendConfirmationError(AccountAlreadyVerified)
)

// @Summary Resend confirmation email
// @Description Send again the confirmation email
// @id ResendConfirmationEmail
// @tags Auth
// @Accept json
// @Produce json
// @Success 202 "Email was sent"
// @Failure 400 {object} ResendConfirmationError
// @Router /users/confirm/resend [post]
// @Param data body users.UserCredentials true "User informations"
func ResendConfirmation(ctx *gin.Context, db *edgedb.Client) {
	var creds users.UserCredentials
	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			ResendInvalidCredentials,
		)
		return
	}
	user, err := creds.Authenticate(db)
	if err != nil && err.Reason == users.InvalidCredentials {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			ResendInvalidCredentials,
		)
		return
	}
	if user.Verified {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			ResendAlreadyVerified,
		)
		return
	}

	originURL := parseRequestOrigin(ctx)
	if err := user.SendConfirmationEmail(originURL); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusAccepted)
}

type EmailConfirmationError string // @name EmailConfirmationError
const (
	AlreadyVerified EmailConfirmationError = EmailConfirmationError(AccountAlreadyVerified)
	InvalidToken    EmailConfirmationError = "InvalidToken"
)

// @Summary Email confirmation
// @Description Confirms a user email using a token
// @id EmailConfirmation
// @tags Auth
// @Accept json
// @Produce json
// @Success 202 "Email was confirmed and account activated"
// @Failure 400 {object} EmailConfirmationError
// @Failure 500 "Server error"
// @Router /users/confirm [get]
// @Param token query string true "Confirmation token"
func ConfirmEmail(ctx *gin.Context, db *edgedb.Client) {
	token := ctx.Query("token")
	logrus.Infof("Received email confirmation token: %s", token)

	userID, tokenValid := users.ValidateEmailConfirmationToken(users.Token(token))
	if !tokenValid {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, InvalidToken)
		return
	}

	user, err := users.FindID(db, userID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if user.Verified {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, AccountAlreadyVerified)
		return
	}

	if err = user.SetActive(true); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	startUserSession(ctx, user.ID)
	ctx.Status(http.StatusAccepted)
}

// @Summary Verify that a password token is valid
// @Description
// @id ValidatePasswordToken
// @tags Auth
// @Success 200 "Password token is valid"
// @Failure 400 "Invalid or expired password reset token"
// @Router /users/password-reset/{token} [get]
// @Param token path string true "Password reset token"
func ValidatePasswordToken(ctx *gin.Context) {
	token := ctx.Param("token")
	_, tokenValid := users.ValidatePasswordResetToken(users.Token(token))
	if !tokenValid {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid token"))
		return
	}
	ctx.Status(http.StatusAccepted)
}

// @Summary Set account password
// @Description Sets a new password for the currently authenticated user
// @id SetPassword
// @tags Auth
// @Success 202 "New password was set"
// @Failure 400 "Invalid password inputs"
// @Failure 403 "Not authenticated"
// @Failure 500 "Database or server error"
// @Router /users/password [post]
// @Param password body users.PasswordInput true "New password"
func SetPassword(ctx *gin.Context, db *edgedb.Client) {
	var newPwd users.PasswordInput
	if err := ctx.BindJSON(&newPwd); err != nil {
		return
	}
	current_user, ok := ctx.Get("current_user")
	if !ok {
		ctx.AbortWithError(
			http.StatusInternalServerError,
			errors.New("failed to retrieve a currently authenticated user"),
		)
		return
	}

	pwd_strength := config.Get().Accounts.PasswordStrength
	newPwd.ValidateStrength(pwd_strength, current_user.(*users.User).InnerUserInput())

	hashed_password, err := newPwd.Hash()
	if err != nil {
		ctx.Error(err)
		return
	}
	query := `update people::User filter .id = current_user_id set {
		password := <str>$0
	}`
	err = db.Execute(context.Background(), query, hashed_password)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusOK)
}

// @Summary Reset account password
// @Description Resets a user's password using a token sent to their email address.
// @id ResetPassword
// @tags Auth
// @Success 202 "Password was reset successfully"
// @Failure 400 "Invalid or expired confirmation token, or invalid input password"
// @Failure 500 "Database error"
// @Router /users/password-reset/{token} [post]
// @Param token path string true "Password reset token"
// @Param password body users.PasswordInput true "New password"
func ResetPassword(ctx *gin.Context) {
	token := ctx.Param("token")
	userID, tokenValid := users.ValidatePasswordResetToken(users.Token(token))
	if !tokenValid {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid token"))
		return
	}

	var newPwd users.PasswordInput

	if err := ctx.BindJSON(&newPwd); err != nil {
		ctx.Error(err)
		return
	}

	pwd_strength := config.Get().Accounts.PasswordStrength
	user, _ := users.FindID(models.DB(), userID)
	if err := newPwd.ValidateStrength(pwd_strength, user.InnerUserInput()); err != nil {
		ctx.Error(err)
		return
	}

	if err := users.SetPassword(userID, &newPwd); err != nil {
		ctx.Error(err)
		return
	}
}

type EmailInput struct {
	Email string `json:"email" binding:"required,email" format:"email"`
} // @name EmailInput

// @Summary Request a password reset token
// @Description A token to reset the password associated to the provided email address is sent, unless the address is not known in the DB.
// @id RequestPasswordReset
// @tags Auth
// @Router /users/forgotten-password [post]
// @Success 202 "Email address is valid and a password reset token was sent"
// @Failure 400 "Invalid email address"
// @Param email body EmailInput true "The email address the account was registered with"
func RequestPasswordReset(ctx *gin.Context, db *edgedb.Client) {
	var email EmailInput
	if err := ctx.BindJSON(&email); err != nil {
		ctx.Error(err)
		return
	}
	user, err := users.Find(db, email.Email)
	if err != nil {
		ctx.Error(err)
		return
	}
	if err = user.RequestPasswordReset(parseRequestOrigin(ctx)); err != nil {
		ctx.Error(err)
		return
	}
}
