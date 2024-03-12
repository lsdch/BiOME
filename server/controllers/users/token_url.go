package accounts

import (
	"darco/proto/config"
	"darco/proto/models/users"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TokenTargetPath struct {
	// The target path when the request origin matches then configured client
	Client string
	// An API endpoint path otherwise
	Default string
}

var emailConfirmationTokenPath = TokenTargetPath{
	Client:  "/email-confirmation",
	Default: "/users/confirm",
}

var passwordResetTokenPath = TokenTargetPath{
	Client:  "/password-reset",
	Default: "/users/password-reset",
}

func newTokenURL(ctx *gin.Context, path TokenTargetPath) users.TokenURL {
	config := config.Get()
	target := users.TokenURL{URL: config.MakeClientURL(path.Client)}

	origin := ctx.Request.Header.Get("Origin")
	originURL, err := url.Parse(origin)

	if err != nil || originURL.Host != target.Host {
		logrus.Infof("Origin URL '%s' does not match configured client URL '%s'. Defaulting to creating token URL pointing to API endpoint.",
			origin, config.Client.DomainName)
		target = users.TokenURL{URL: config.MakeURL(path.Default)}
	}
	return target
}
