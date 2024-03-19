package accounts

import (
	"darco/proto/utils"
	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	CONFIRM_EMAIL_PATH  = "/api/v1/users/confirm"
	PASSWORD_RESET_PATH = "/api/v1/users/password-reset"
)

func newRedirectURL(ctx *gin.Context, path string) utils.RedirectURL {
	target := utils.RedirectURL{
		URL: url.URL{
			Scheme: ctx.Request.URL.Scheme,
			Host:   ctx.Request.Host,
			Path:   path,
		},
	}
	if redirect := ctx.Query("redirect"); redirect != "" {
		target.SetRedirect(redirect)
	}
	return target
}

func confirmEmailURL(ctx *gin.Context) utils.RedirectURL {
	return newRedirectURL(ctx, CONFIRM_EMAIL_PATH)
}

func passwordResetURL(ctx *gin.Context) utils.RedirectURL {
	return newRedirectURL(ctx, PASSWORD_RESET_PATH)
}
