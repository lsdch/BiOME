package resolvers

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

type OwnershipResolver interface {
	IsOwner() bool
}

type OwnerOr[R AccessResolver, O OwnershipResolver] struct {
	AccessResolver R
	Ownership      O
}

func (i *OwnerOr[R, O]) Resolve(ctx huma.Context) []error {
	i.AccessResolver.ResolveAuth(ctx)
	if i.AccessResolver.IsGranted() || i.Ownership.IsOwner() {
		return nil
	} else {
		return []error{
			huma.Error401Unauthorized(
				fmt.Sprintf(
					"Access restricted to the ressource owner/maintainer or %s users",
					i.AccessResolver.RoleRequired(),
				),
			),
		}
	}
}
