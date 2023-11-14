package middlewares

import (
	"darco/proto/config"
	"darco/proto/models"
	"darco/proto/models/users"
	"darco/proto/services/tokens"
	"strings"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func AuthenticationMiddleware(ctx *gin.Context) {

	var access_token string
	cookie, err := ctx.Cookie("token")

	ctx.Set("currentUser", nil)
	ctx.Set("db", models.DB())

	authorizationHeader := ctx.Request.Header.Get("Authorization")
	fields := strings.Fields((authorizationHeader))
	if len(fields) != 0 && fields[0] == "Bearer" {
		access_token = fields[1]
	} else if err == nil {
		access_token = cookie
	}

	if access_token == "" {
		logrus.Debugf("No authentication token")
		return
	}

	sub, err := tokens.ValidateToken(config.Get(), access_token)
	if err != nil {
		logrus.Debugf("Invalid token received")
		return
	}

	userID, err := uuid.Parse(sub.(string))
	if err != nil {
		logrus.Debugf("Token %s does not hold a valid UUID", sub)
		return
	}

	current_user, err := users.FindID(models.DB(), userID)
	if err != nil {
		logrus.Errorf("Token was validated but does not match an existing user.")
		return
	}

	logrus.Debugf("User authenticated %+v", current_user)

	client := models.DB().WithGlobals(map[string]interface{}{"current_user_id": edgedb.UUID(userID)})
	ctx.Set("current_user", current_user)
	ctx.Set("db", client)
	ctx.Next()
}
