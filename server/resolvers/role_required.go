package resolvers

import (
	"darco/proto/models/people"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

type RoleSpecifier interface {
	Role() people.UserRole
}

type Contributor struct{}

func (r Contributor) Role() people.UserRole {
	return people.Contributor
}

var _ RoleSpecifier = (*Contributor)(nil)

type Maintainer struct{}

func (r Maintainer) Role() people.UserRole {
	return people.Maintainer
}

var _ RoleSpecifier = (*Maintainer)(nil)

type Admin struct{}

func (r Admin) Role() people.UserRole {
	return people.Admin
}

var _ RoleSpecifier = (*Admin)(nil)

type AccessResolver interface {
	UserResolver
	IsGranted() bool
	RoleRequired() people.UserRole
}

type AccessRestricted[R RoleSpecifier] struct {
	AuthResolver
	RoleSpec R
}

func (a *AccessRestricted[R]) IsGranted() bool {
	return a.AuthUser().Role.IsGreaterEqual(a.RoleSpec.Role())
}

func (a *AccessRestricted[RV]) Resolve(ctx huma.Context) []error {
	a.AuthResolver.ResolveAuth(ctx)
	if !a.IsGranted() {
		return []error{
			huma.Error401Unauthorized(
				fmt.Sprintf("Access restricted to %s users", a.RoleSpec.Role()),
			),
		}
	}
	return nil
}

func (a *AccessRestricted[R]) RoleRequired() people.UserRole {
	return a.RoleSpec.Role()
}

var _ AccessResolver = (*AccessRestricted[Contributor])(nil)

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
					"Access restricted to the ressource creator or %s users",
					i.AccessResolver.RoleRequired(),
				),
			),
		}
	}
}
