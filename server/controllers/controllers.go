package controllers

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

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
	ID edgedb.UUID `path:"id" format:"uuid"`
}

func (i UUIDInput) Identifier() edgedb.UUID {
	return i.ID
}

type CodeInput struct {
	Code string `path:"code"`
}

func (i CodeInput) Identifier() string {
	return i.Code
}

type SlugInput struct {
	Slug string `path:"slug"`
}

func (i SlugInput) Identifier() string {
	return i.Slug
}

// Implementation assertions
var _ IdentifierInput[edgedb.UUID] = (*UUIDInput)(nil)
var _ IdentifierInput[string] = (*CodeInput)(nil)

// A simple response output that carries a message
type Message struct {
	Body string
}
