package accounts

import (
	"darco/proto/config"
	"darco/proto/models/users"
	_ "darco/proto/models/validations"
	"darco/proto/services/tokens"
	"net/http"
	"net/url"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary Delete an account
// @Description Deletes an account
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
