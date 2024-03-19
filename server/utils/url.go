package utils

import "net/url"

type URLQuery interface {
	Params() url.Values
	Encode() *url.URL
	String() string
}

// A URL that may embed a '?redirect=' query parameter.
type RedirectURL struct {
	url.URL
	redirect string `faker:"-"`
}

func (u *RedirectURL) SetRedirect(path string) {
	u.redirect = path
}
func (u *RedirectURL) Redirect() string {
	return u.redirect
}

func (u *RedirectURL) String() string {
	return u.Encode().String()
}

func (u *RedirectURL) Params() url.Values {
	params := url.Values{}
	if u.redirect != "" {
		params.Set("redirect", u.redirect)
	}
	return params
}

func (u *RedirectURL) Encode() *url.URL {
	u.URL.RawQuery = u.Params().Encode()
	return &u.URL
}
