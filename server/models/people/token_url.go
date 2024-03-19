package people

import (
	"darco/proto/utils"
	"net/url"
)

// A URL that embeds '?token=...&redirect=...' query parameters.
// Redirect is optional.
type TokenURL struct {
	utils.RedirectURL
	token Token `faker:"-"`
}

func NewTokenURL(url utils.RedirectURL, token Token) TokenURL {
	return TokenURL{url, token}
}

func (u *TokenURL) SetToken(token Token) {
	u.token = token
}
func (u *TokenURL) Token() Token {
	return u.token
}

func (u *TokenURL) Params() url.Values {
	params := url.Values{}
	if u.token != "" {
		params.Set("token", string(u.token))
	}
	for k, values := range u.RedirectURL.Params() {
		for _, v := range values {
			params.Set(k, v)
		}
	}
	return params
}

func (u *TokenURL) Encode() *url.URL {
	u.URL.RawQuery = u.Params().Encode()
	return &u.URL
}

func (u *TokenURL) String() string {
	return u.Encode().String()
}
