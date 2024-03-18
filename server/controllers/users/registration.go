package accounts

import (
	"darco/proto/models/people"
	users "darco/proto/models/people"
	_ "darco/proto/models/validations"
	"fmt"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary Register user
// @Description Register a new user account, that is inactive (until email is verified or admin intervention), and has role 'Guest'
// @id RegisterUser
// @tags Auth
// @Accept json
// @Produce json
// @Success 202 "User created and waiting for email verification"
// @Failure 400 {object} validations.FieldErrors
// @Router /users/register [post]
// @Param data body people.PendingUserRequestInput true "User informations"
// @Param redirect query string false "Path to redirect to when email is successfully confirmed"
func Register(ctx *gin.Context, db *edgedb.Client) {
	var newUser people.PendingUserRequestInput
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.Error(err)
		return
	}

	logrus.Infof("Attempting to create account for %s %s (%s)",
		newUser.Person.FirstName,
		newUser.Person.LastName,
		newUser.User.Email,
	)

	pending, err := newUser.Register(db)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	target := confirmEmailURL(ctx)
	if err := pending.User.SendConfirmationEmail(db, target); err != nil {
		msg := fmt.Errorf(
			"Failed to send account confirmation email to '%s'. This is likely to be a server configuration error.",
			newUser.User.Email,
		)
		logrus.Errorf("%s Error: %v", msg, err)
		ctx.AbortWithError(http.StatusInternalServerError, msg)
		return
	}
	ctx.Status(http.StatusAccepted)
}

type EmailConfirmationError string // @name EmailConfirmationError
const (
	AlreadyVerified EmailConfirmationError = "AlreadyVerified"
	InvalidToken    EmailConfirmationError = "InvalidToken"
)

// @Summary Email confirmation
// @Description Confirms a user email using a token
// @id EmailConfirmation
// @tags Auth
// @Accept json
// @Produce json
// @Success 200 "Email was confirmed and account activated"
// @Failure 400 {object} EmailConfirmationError
// @Failure 500 "Server error"
// @Router /users/confirm [get]
// @Param token query string true "Confirmation token"
// @Param redirect query string false "Path to redirect to on success"
func ConfirmEmail(ctx *gin.Context, db *edgedb.Client) {
	token := ctx.Query("token")
	logrus.Infof("Received email confirmation token: %s", token)

	user, tokenValid := users.ValidateAccountToken(db, users.Token(token), users.EmailConfirmationToken)
	if !tokenValid {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, InvalidToken)
		return
	}

	if user.EmailConfirmed {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, AlreadyVerified)
		return
	}

	if err := user.SetEmailConfirmed(db, true); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if user.IsActive {
		startUserSession(ctx, user)
	}
	if redirect := ctx.Param("redirect"); redirect != "" {
		ctx.Redirect(http.StatusOK, redirect)
	} else {
		ctx.Status(http.StatusOK)
	}
}

type ResendConfirmationError string // @name ResendConfirmationError

const (
	ResendInvalidCredentials ResendConfirmationError = "InvalidCredentials"
	ResendAlreadyVerified    ResendConfirmationError = ResendConfirmationError(AlreadyVerified)
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
// @Param redirect query string false "Redirect to path on confirmation"
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
	if user.EmailConfirmed {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			ResendAlreadyVerified,
		)
		return
	}

	tokenURL := confirmEmailURL(ctx)

	if err := user.SendConfirmationEmail(db, tokenURL); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusAccepted)
}
