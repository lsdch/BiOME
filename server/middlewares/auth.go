package middlewares

import (
	"darco/proto/config"
	"darco/proto/db"
	"darco/proto/models/users"
	"darco/proto/services/tokens"
	"strings"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Gin context key where the current user is stored as a pointer, if any
const CTX_CURRENT_USER_KEY = "current_user"

// Gin context key where the database client is stored as a pointer
const CTX_DATABASE_KEY = "db"

func AuthenticationMiddleware(ctx *gin.Context) {

	var access_token string
	cookie, err := ctx.Cookie("token")

	ctx.Set(CTX_CURRENT_USER_KEY, nil)
	ctx.Set(CTX_DATABASE_KEY, db.Client())

	authorizationHeader := ctx.Request.Header.Get("Authorization")
	fields := strings.Fields((authorizationHeader))
	if len(fields) != 0 && fields[0] == "Bearer" {
		access_token = fields[1]
	} else if err == nil {
		access_token = cookie
	}

	if access_token == "" {
		logrus.Debugf("Auth middleware: No authentication token")
		return
	}

	sub, err := tokens.ValidateToken(config.Get(), access_token)
	if err != nil {
		logrus.Debugf("Auth middleware: Invalid token received")
		return
	}

	userID, err := edgedb.ParseUUID(sub.(string))
	if err != nil {
		logrus.Debugf("Auth middleware: Token %s does not hold a valid UUID", sub)
		return
	}

	current_user, err := users.FindID(db.Client(), userID)
	if err != nil {
		logrus.Errorf("Auth middleware: Token was validated but does not match an existing user.")
		return
	}

	// Authentication succeeded
	logrus.Debugf("Auth middleware: User authenticated %+v", current_user)
	client := db.WithCurrentUser(userID)
	ctx.Set(CTX_CURRENT_USER_KEY, current_user)
	ctx.Set(CTX_DATABASE_KEY, client)
	ctx.Next()
}
