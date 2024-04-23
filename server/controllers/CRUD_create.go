package controllers

import (
	"context"
	"darco/proto/models"
	"darco/proto/resolvers"

	"github.com/danielgtaylor/huma/v2"
)

type CreateHandlerInput[Item models.Creatable[Created], Created any] struct {
	resolvers.AuthRequired
	Body Item
}

type CreateHandlerOutput[Created any] struct {
	Body Created
}

func CreateHandler[
	Input models.Creatable[Created],
	Created any,
](ctx context.Context, input *CreateHandlerInput[Input, Created]) (*CreateHandlerOutput[Created], error) {
	created, err := input.Body.Create(input.DB())
	if err != nil {
		return nil, huma.Error500InternalServerError("Item creation failed", err)
	}
	return &CreateHandlerOutput[Created]{Body: created}, nil
}
