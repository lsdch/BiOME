package accounts

import (
	"context"
	"darco/proto/config"
	"darco/proto/models/users"
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
	_, tokenValid := users.ValidatePasswordResetToken(db, users.Token(token))
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
	var newPwd users.PasswordInput
	if err := ctx.BindJSON(&newPwd); err != nil {
		return
	}

	pwd_strength := config.Get().Accounts.PasswordStrength
	err := newPwd.ValidateStrength(pwd_strength, user.InnerUserInput())
	if err != nil {
		ctx.Error(err)
		return
	}

	query := `update people::User filter .id = current_user_id set { password := <str>$0 }`
	err = db.Execute(context.Background(), query, newPwd)
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
func ResetPassword(ctx *gin.Context, db *edgedb.Client) {
	token := ctx.Param("token")
	userID, tokenValid := users.ValidatePasswordResetToken(db, users.Token(token))
	if !tokenValid {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid token"))
		return
	}

	// User must exist if token validation succeeded
	user, _ := users.FindID(db, userID)

	var newPwd users.PasswordInput
	if err := ctx.BindJSON(&newPwd); err != nil {
		ctx.Error(err)
		return
	}

	pwd_strength := config.Get().Accounts.PasswordStrength
	if err := newPwd.ValidateStrength(pwd_strength, user.InnerUserInput()); err != nil {
		ctx.Error(err)
		return
	}

	if err := user.SetPassword(db, newPwd); err != nil {
		ctx.Error(err)
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
	if err = user.RequestPasswordReset(parseRequestOrigin(ctx)); err != nil {
		ctx.Error(err)
		return
	}
}
