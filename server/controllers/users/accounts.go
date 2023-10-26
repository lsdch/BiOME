package accounts

import (
	"darco/proto/config"
	"darco/proto/models/users"
	"darco/proto/services/tokens"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
	Token string `json:"token" example:"some-generated-jwt"`
} // @name TokenResponse

// @Summary Authenticate user
// @Description Authenticate user with their credentials and set a JWT.
// @id Login
// @tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} TokenResponse "Returns a token and stores it as a session cookie"
// @Failure 400 "Invalid credentials"
// @Router /login [post]
// @Param data body users.UserCredentials true "User credentials"
func Login(ctx *gin.Context) {
	var credentials users.UserCredentials
	err := ctx.ShouldBindJSON(&credentials)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := credentials.Authenticate()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	config := config.Get()
	token, err := tokens.GenerateToken(user.ID, config.TokenLifetime)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	lifetime := -1
	if credentials.Remember {
		lifetime = int(config.TokenLifetime.Seconds())
	}
	ctx.SetCookie("token", token, lifetime, "/", "", true, true)
	ctx.JSON(http.StatusOK, TokenResponse{token})
}

// @Summary Register user
// @Description Register a new user account, that is inactive (until email is verified or admin intervention), and has role 'Guest'
// @id RegisterUser
// @tags Auth
// @Accept json
// @Produce json
// @Success 202 "User created and waiting for email verification"
// @Failure 400 "Invalid parameters"
// @Router /users/register [post]
// @Param data body users.UserInput true "User informations"
func Register(ctx *gin.Context) {
	var newUser users.UserInput
	err := ctx.ShouldBindJSON(&newUser)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err = newUser.Register(config.Get()); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.Status(http.StatusAccepted)
}

// @Summary Resend confirmation email
// @Description Send again the confirmation email
// @id ResendConfirmationEmail
// @tags Auth
// @Accept json
// @Produce json
// @Success 202 "Email was sent"
// @Failure 400 "Invalid parameters"
// @Router /users/confirm/resend [post]
// @Param data body users.UserCredentials true "User informations"
func ResendConfirmation(ctx *gin.Context) {
	var creds users.UserCredentials
	err := ctx.ShouldBindJSON(&creds)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := users.Find(creds.Identifier)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if user.Verified {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("user email already verified"))
		return
	}
	err = user.SendConfirmationEmail()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusAccepted)
}

// @Summary Email confirmation
// @Description Confirms a user email using a token
// @id EmailConfirmation
// @tags Auth
// @Accept json
// @Produce json
// @Success 202 "Email was confirmed and account activated"
// @Failure 400 "Invalid or expired confirmation token"
// @Failure 500 "Token parse error"
// @Router /users/confirm [get]
// @Param token query string true "Confirmation token"
func ConfirmEmail(ctx *gin.Context) {
	token := ctx.Query("token")

	email, err := tokens.ValidateToken(config.Get(), token)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := users.Find(email.(string))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if user.Verified {
		err := errors.New("this account was already verified")
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = user.SetActive(true)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

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
	_, tokenValid := users.ValidatePasswordResetToken(token)
	if !tokenValid {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid token"))
		return
	}
	ctx.Status(http.StatusAccepted)
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
	userID, tokenValid := users.ValidatePasswordResetToken(token)
	if !tokenValid {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid token"))
		return
	}

	var newPwd users.PasswordInput
	if err := ctx.BindJSON(&newPwd); err != nil {
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
func RequestPasswordReset(ctx *gin.Context) {
	var email EmailInput
	if err := ctx.BindJSON(&email); err != nil {
		return
	}
	user, err := users.Find(email.Email)
	if err != nil {
		ctx.Error(err)
		return
	}
	if err = user.RequestPasswordReset(); err != nil {
		ctx.Error(err)
		return
	}
}
