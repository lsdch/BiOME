package controllers

import (
	"context"
	"darco/proto/models"
	"darco/proto/resolvers"
)

type CreateInputBody[Item models.Creatable[Created], Created any] interface {
	resolvers.AuthDBProvider
	Item() Item
}

type CreateHandlerInput[Item models.Creatable[Created], Created any] struct {
	resolvers.AuthRequired
	Body Item
}

func (i CreateHandlerInput[Item, Created]) Item() Item {
	return i.Body
}

type CreateHandlerOutput[Created any] struct {
	Body Created
}

func CreateHandler[
	Item models.Creatable[Created],
	Created any,
](ctx context.Context, input *CreateHandlerInput[Item, Created]) (*CreateHandlerOutput[Created], error) {
	return CreateHandlerWithInput(ctx, input)
}

func CreateHandlerWithInput[
	Input CreateInputBody[Item, Created],
	Item models.Creatable[Created],
	Created any,
](ctx context.Context, input Input) (*CreateHandlerOutput[Created], error) {
	created, err := input.Item().Create(input.DB())
	if err = StatusError(err); err != nil {
		return nil, err
	}
	return &CreateHandlerOutput[Created]{Body: created}, nil
}
