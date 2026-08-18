package auth

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/models"
)

type Policy interface {
	Eval(ctx *models.User) bool
	IsAuthRequired() bool
}

type PublicPolicy struct{}

func (p PublicPolicy) Eval(user *models.User) bool {
	return true
}

func (p PublicPolicy) IsAuthRequired() bool {
	return false
}

type AuthenticatedPolicy struct{}

func (p AuthenticatedPolicy) Eval(user *models.User) bool {
	return user != nil
}

func (p AuthenticatedPolicy) IsAuthRequired() bool {
	return true
}

type RolePolicy struct {
	Role models.UserRole `json:"role"`
}

func (p RolePolicy) Eval(user *models.User) bool {
	if user == nil {
		return false
	}
	return user.Role.IsGreaterEqual(p.Role)
}

func (p RolePolicy) IsAuthRequired() bool {
	return true
}

type AndPolicy struct {
	Policies []Policy `json:"all"`
}

func (p AndPolicy) Eval(user *models.User) bool {
	for _, pol := range p.Policies {
		if !pol.Eval(user) {
			return false
		}
	}
	return true
}

func (p AndPolicy) IsAuthRequired() bool {
	for _, pol := range p.Policies {
		if pol.IsAuthRequired() {
			return true
		}
	}
	return false
}

type OrPolicy struct {
	Policies []Policy `json:"any"`
}

func (p OrPolicy) Eval(user *models.User) bool {
	for _, pol := range p.Policies {
		if pol.Eval(user) {
			return true
		}
	}
	return false
}

func (p OrPolicy) IsAuthRequired() bool {
	for _, pol := range p.Policies {
		// Authentication is required if ALL policies require authentication. If any policy does not require authentication, then the OR policy does not require authentication.
		if !pol.IsAuthRequired() {
			return false
		}
	}
	return true
}

type NotPolicy struct {
	Policy Policy
}

func (p NotPolicy) Eval(user *models.User) bool {
	return !p.Policy.Eval(user)
}

func (p NotPolicy) IsAuthRequired() bool {
	return p.Policy.IsAuthRequired()
}

func Public() Policy {
	return PublicPolicy{}
}

func Authenticated() Policy {
	return AuthenticatedPolicy{}
}

func Role(role models.UserRole) Policy {
	return RolePolicy{Role: role}
}

func And(policies ...Policy) Policy {
	return AndPolicy{Policies: policies}
}

func Or(policies ...Policy) Policy {
	return OrPolicy{Policies: policies}
}

func Not(policy Policy) Policy {
	return NotPolicy{Policy: policy}
}

func WithPolicy(op huma.Operation, p Policy) huma.Operation {
	if op.Extensions == nil {
		op.Extensions = map[string]any{}
	}
	op.Extensions["x-policy"] = p
	return op
}
