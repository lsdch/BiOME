package resolvers

import (
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirupsen/logrus"
)

type HostResolver struct {
	Host   string
	URL    url.URL
	Origin url.URL
}

func (i *HostResolver) Resolve(ctx huma.Context) []error {
	i.URL = ctx.URL()
	i.Host = ctx.Host()
	originURL, err := url.Parse(ctx.Header("Origin"))
	if err != nil {
		logrus.Errorf("[HostResolver] Failed to parse request origin: %v", err)
	} else {
		i.Origin = *originURL
	}
	return nil
}

// Generates a URL with the provided path on the request origin (e.g. the client domain)
func (i HostResolver) OriginPath(path string) url.URL {
	return url.URL{
		Host:   i.Origin.Host,
		Scheme: "https",
		Path:   path,
	}
}

// Generates a URL for the provided path on the API backend
func (i *HostResolver) GenerateURL(path string) url.URL {
	return url.URL{
		Host:   i.Host,
		Scheme: "http",
		Path:   path,
	}
}
