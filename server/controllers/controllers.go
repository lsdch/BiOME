package controllers

import (
	"reflect"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirupsen/logrus"
)

type ResponseBody[T any] struct {
	Body T
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
	ID geltypes.UUID `path:"id" format:"uuid"`
}

func (i UUIDInput) Identifier() geltypes.UUID {
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

type LabelInput struct {
	Label string `path:"label"`
}

func (i LabelInput) Identifier() string {
	return i.Label
}

type EmailInput struct {
	Email string `path:"email"`
}

func (i EmailInput) Identifier() string {
	return i.Email
}

// Implementation assertions
var _ IdentifierInput[geltypes.UUID] = (*UUIDInput)(nil)
var _ IdentifierInput[string] = (*CodeInput)(nil)

// A simple response output that carries a message
type Message struct {
	Body string
}

// StatusError transforms errors resulting from DB calls into Huma status errors
func StatusError(err error) huma.StatusError {
	if err == nil {
		return nil
	}
	if e, ok := err.(huma.StatusError); ok {
		return e
	}
	if db.IsNoData(err) {
		return huma.Error404NotFound("Item not found", err)
	}
	if isConstraintErr, constraintErr := db.IsConstraintViolation(err); isConstraintErr {
		return huma.Error422UnprocessableEntity("Invalid input", constraintErr)
	}

	// Other errors are HTTP 500
	logrus.Errorf("Server error: %v", err)
	return huma.Error500InternalServerError("Server error", err)
}
