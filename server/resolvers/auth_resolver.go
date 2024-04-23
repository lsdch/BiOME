package resolvers

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/services/tokens"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

// A resolver that provides a DB client, having set the currently authenticated user (if any)
type AuthDBProvider interface {
	DB() *edgedb.Client
}

type AuthResolver struct {
	*people.User
	AuthToken       string      // Auth token parsed from session cookie or authorization header
	AuthTokenHeader string      `header:"Authorization" doc:"Authorization header formatted as \"Bearer auth_token\". Takes precedence over session cookie if set." example:"Bearer <JWT string>"`
	Session         http.Cookie `cookie:"auth_token" doc:"Session cookie containing JWT"`
}

func (p *AuthResolver) DB() *edgedb.Client {
	if p.User != nil {
		return db.WithCurrentUser(p.ID)
	} else {
		return db.Client()
	}
}

func (p *AuthResolver) Resolve(ctx huma.Context) []error {
	p.User = nil
	var accessToken string

	logrus.Debugf("Session cookie: %+v", p.Session)

	if p.AuthTokenHeader != "" {
		accessToken = strings.TrimPrefix(p.AuthTokenHeader, "Bearer ")
	} else {
		accessToken = p.Session.Value
	}

	if accessToken == "" {
		logrus.Debugf("Auth middleware: No authentication token")
		return nil
	}

	sub, err := tokens.ValidateToken(accessToken)
	if err != nil {
		logrus.Debugf("Auth middleware: Invalid token received %v", err)
		return nil
	}

	userID, err := edgedb.ParseUUID(sub.(string))
	if err != nil {
		logrus.Debugf("Auth middleware: Token %s does not hold a valid UUID", sub)
		return nil
	}

	currentUser, err := people.FindID(db.Client(), userID)
	if err != nil {
		logrus.Errorf("Auth middleware: Token was validated but does not match an existing user.")
		return nil
	}

	logrus.Debugf("Auth middleware: User authenticated %+v", currentUser)
	p.AuthToken = accessToken
	p.User = currentUser
	return nil
}
