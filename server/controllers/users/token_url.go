package accounts

import (
	users "darco/proto/models/people"
	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	CONFIRM_EMAIL_PATH  = "/api/v1/users/confirm"
	PASSWORD_RESET_PATH = "/api/v1/users/password-reset"
)

func newTokenURL(ctx *gin.Context, path string) users.TokenURL {
	target := users.NewTokenURL(url.URL{
		Scheme: ctx.Request.URL.Scheme,
		Host:   ctx.Request.Host,
		Path:   path,
	})
	if redirect := ctx.Query("redirect"); redirect != "" {
		target.SetRedirect(redirect)
	}
	return target
}

func confirmEmailURL(ctx *gin.Context) users.TokenURL {
	return newTokenURL(ctx, CONFIRM_EMAIL_PATH)
}

func passwordResetURL(ctx *gin.Context) users.TokenURL {
	return newTokenURL(ctx, PASSWORD_RESET_PATH)
}
