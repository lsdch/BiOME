package accounts

import (
	users "darco/proto/models/people"
	"darco/proto/models/settings"
	_ "darco/proto/models/validations"
	"darco/proto/services/tokens"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

type TokenResponse struct {
	Token string `edgedb:"token" json:"token"`
}

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
	token := startUserSession(ctx, user)
	ctx.JSON(http.StatusOK, TokenResponse{token})
}

// Stores authentication token in cookies after a user has logged, and send it back with HTTP status OK.
func startUserSession(ctx *gin.Context, user *users.User) string {
	token, err := user.GenerateJWT()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return ""
	}

	ctx.SetCookie(
		tokens.AUTH_TOKEN_COOKIE,
		token,
		int(settings.Security().AuthTokenLifetimeSeconds),
		"/",
		ctx.Request.Host,
		true,
		true,
	)
	return token
}

func unsetCookie(ctx *gin.Context, name string) {
	ctx.SetCookie(name, "", -1, "/", ctx.Request.Host, false, true)
}

// @Summary Logout user
// @Description Log out currently authenticated user by revoking authentication token in cookies
// @id Logout
// @tags Auth
// @Accept json
// @Produce json
// @Success 200 "User logged out"
// @Router /logout [post]
func Logout(ctx *gin.Context) {
	unsetCookie(ctx, tokens.AUTH_TOKEN_COOKIE)
	unsetCookie(ctx, tokens.REFRESH_TOKEN_COOKIE)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
