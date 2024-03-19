package people_test

import (
	"darco/proto/models/people"
	"darco/proto/utils"
	"fmt"
	"net/url"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func TestTokenURL(t *testing.T) {
	fakeURL := utils.RedirectURL{
		URL: url.URL{
			Host:   faker.DomainName(),
			Path:   "/some/path",
			Scheme: "https",
		},
	}

	token := people.Token(faker.Password())
	tokenURL := people.NewTokenURL(fakeURL, token)
	assert.Equal(t, token, tokenURL.Token())
	t.Run("TokenURL set token", func(t *testing.T) {
		newToken := people.Token(faker.Password())
		tokenURL.SetToken(newToken)
		assert.Equal(t, newToken, tokenURL.Token())
		fakeURL.RawQuery = fmt.Sprintf("token=%s", newToken)
		assert.Equal(t, fakeURL.URL.String(), tokenURL.String())
	})

}
