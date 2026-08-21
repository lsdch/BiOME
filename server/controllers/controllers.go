package controllers

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/types"

	"github.com/danielgtaylor/huma/v2"
)

type Controller interface {
	RegisterRoutes(r *router.Router)
}

type BodyTransporter[Item any] struct {
	Body Item
}

type IdentifierInput[T any] interface {
	Identifier() T
}

type StrIdentifier string

func (i StrIdentifier) Identifier() string {
	return string(i)
}

func (o StrIdentifier) Schema(r huma.Registry) *huma.Schema {
	return r.Schema(reflect.TypeOf(""), true, "")
}

type UUIDInput struct {
	ID uuid.UUID `path:"id" format:"uuid"`
}

func (i UUIDInput) Identifier() uuid.UUID {
	return i.ID
}

type ULIDPath struct {
	ULID types.ULID `path:"ulid" format:"ulid"`
}

type CodePath struct {
	Code string `path:"code"`
}

func (i CodePath) Identifier() string {
	return i.Code
}

type SlugInput struct {
	Slug string `path:"slug"`
}

func (i SlugInput) Identifier() string {
	return i.Slug
}

type LabelInput struct {
	Label string `path:"label"`
}

func (i LabelInput) Identifier() string {
	return i.Label
}

type NameInput struct {
	Name string `path:"name"`
}

func (i NameInput) Identifier() string {
	return i.Name
}

type EmailInput struct {
	Email string `path:"email"`
}

func (i EmailInput) Identifier() string {
	return i.Email
}

type NumberInput struct {
	Number int64 `path:"number"`
}

func (i NumberInput) Identifier() int64 {
	return i.Number
}

// Implementation assertions
var _ IdentifierInput[string] = (*CodePath)(nil)
