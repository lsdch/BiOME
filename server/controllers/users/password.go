package accounts

import (
	"darco/proto/models/users"
	"darco/proto/models/validations"
	"errors"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

// @Summary Verify that a password token is valid
// @Description
// @id ValidatePasswordToken
// @tags Auth
// @Success 200 "Password token is valid"
// @Failure 400 "Invalid or expired password reset token"
// @Router /users/password-reset/{token} [get]
// @Param token path string true "Password reset token"
func ValidatePasswordToken(ctx *gin.Context, db *edgedb.Client) {
	token := ctx.Param("token")
	_, tokenValid := users.ValidateAccountToken(db, users.Token(token), users.PasswordResetToken)
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
// @Router /account/password [post]
// @Param password body users.PasswordInput true "New password"
func SetPassword(ctx *gin.Context, db *edgedb.Client, user *users.User) {
	var newPwd users.NewPasswordInput
	if err := ctx.BindJSON(&newPwd); err != nil {
		return
	}

	if !user.PasswordMatch(db, newPwd.Password) {
		ctx.Error(validations.InputValidationError{
			Field:   "password",
			Message: "Invalid password",
		})
		return
	}

	if err := user.SetPassword(db, newPwd.NewPassword.Password); err != nil {
		ctx.Error(weakPasswordError("new_password.password", err))
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
func ResetPassword(ctx *gin.Context, db *edgedb.Client) {
	token := ctx.Param("token")
	user, tokenValid := users.ValidateAccountToken(db, users.Token(token), users.PasswordResetToken)
	if !tokenValid {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid token"))
		return
	}

	var newPwd users.PasswordInput
	if err := ctx.BindJSON(&newPwd); err != nil {
		ctx.Error(err)
		return
	}

	if err := user.SetPassword(db, newPwd.Password); err != nil {
		ctx.Error(weakPasswordError("password", err))
		return
	}
	ctx.Status(http.StatusAccepted)
}

// Single email input type used for password reset requests
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

	tokenURL := newTokenURL(ctx, passwordResetTokenPath)
	if err = user.RequestPasswordReset(db, tokenURL); err != nil {
		ctx.Error(err)
		return
	}
}
