package resolvers

import (
	"net/http"
	"strings"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/users"
	"github.com/lsdch/biome/services"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirupsen/logrus"
)

// A resolver that provides a DB client, having set the currently authenticated user (if any)
type AuthDBProvider interface {
	DB() *gel.Client
}

type UserResolver interface {
	AuthDBProvider
	AuthUser() (*users.User, bool)
	ResolveAuth(huma.Context)
}

type AuthResolver struct {
	User            *users.User
	AuthToken       string      // Auth token parsed from session cookie or authorization header
	AuthTokenHeader string      `header:"Authorization" doc:"Authorization header formatted as \"Bearer auth_token\". Takes precedence over session cookie if set." example:"Bearer <JWT string>"`
	Session         http.Cookie `cookie:"auth_token" doc:"Session cookie containing JWT"`
}

func (p *AuthResolver) IsGranted(role biomedb.UserRole) bool {
	if p.User == nil {
		return false
	}
	return p.User.IsGranted(role)
}

func (p *AuthResolver) AuthUser() (*users.User, bool) {
	if p.User != nil {
		return p.User, true
	}
	return nil, false
}

func (p *AuthResolver) Resolve(ctx huma.Context) []error {
	p.ResolveAuth(ctx)
	return nil
}

func (p *AuthResolver) ResolveAuth(ctx huma.Context) {
	p.User = nil
	p.AuthToken = ""
	var accessToken string

	if p.AuthTokenHeader != "" {
		accessToken = strings.TrimPrefix(p.AuthTokenHeader, "Bearer ")
	} else {
		accessToken = p.Session.Value
	}

	if accessToken == "" {
		logrus.Debugf("Auth middleware: No authentication token")
		return
	}

	sub, err := services.ValidateJWT(accessToken)
	if err != nil {
		logrus.Debugf("Auth middleware: Invalid token received [%v]", err)
		return
	}

	userID, err := geltypes.ParseUUID(sub.(string))
	if err != nil {
		logrus.Debugf("Auth middleware: Token %s does not hold a valid UUID", sub)
		return
	}

	currentUser, err := people.FindID(db.Client(), userID)
	if err != nil {
		logrus.Errorf("Auth middleware: Token was validated but does not match an existing user.")
		return
	}

	logrus.Debugf(
		"Auth middleware: authenticated %s [%s]",
		currentUser.Person.FullName, currentUser.Role,
	)

	p.AuthToken = accessToken
	p.User = &currentUser
}

var _ UserResolver = (*AuthResolver)(nil)
