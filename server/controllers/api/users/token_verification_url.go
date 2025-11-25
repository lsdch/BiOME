package accounts

import (
	"encoding/json"
	"net/url"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type TokenVerificationURL url.URL

func (u TokenVerificationURL) Schema(r huma.Registry) *huma.Schema {
	s := r.Schema(reflect.TypeFor[url.URL](), false, "")
	s.Description = "A URL used to generate the verification link, which can be set by the web client. Verification token will be added as a URL query parameter."
	return s
}

func (u TokenVerificationURL) URL() url.URL {
	return url.URL(u)
}

func (u *TokenVerificationURL) UnmarshalJSON(data []byte) error {
	var stringURL string
	if err := json.Unmarshal(data, &stringURL); err != nil {
		return err
	}
	parsedURL, err := url.Parse(stringURL)
	if err != nil {
		return err
	}
	*u = TokenVerificationURL(*parsedURL)
	return nil
}
