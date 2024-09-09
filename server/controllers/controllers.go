package controllers

import (
	"github.com/edgedb/edgedb-go"
)

type IdentifierInput[T any] interface {
	Identifier() T
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
