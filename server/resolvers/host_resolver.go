package resolvers

import (
	"net/url"

	"github.com/danielgtaylor/huma/v2"
)

type HostResolver struct {
	Host string
	URL  url.URL
}

func (i *HostResolver) Resolve(ctx huma.Context) []error {
	i.URL = ctx.URL()
	i.Host = ctx.URL().Host
	return nil
}

func (i *HostResolver) GenerateURL(path string) url.URL {
	return url.URL{
		Host:   i.Host,
		Scheme: i.URL.Scheme,
		Path:   path,
	}
}
