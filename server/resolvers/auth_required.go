package resolvers

import "github.com/danielgtaylor/huma/v2"

type AuthRequired struct {
	AuthResolver
}

func (a *AuthRequired) Resolve(ctx huma.Context) []error {
	a.AuthResolver.Resolve(ctx)
	if a.User == nil {
		return []error{huma.Error401Unauthorized("Authentication required")}
	}
	return nil
}
