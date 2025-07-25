package controllers

import (
	"context"

	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/resolvers"
)

type CreateInputBody[Item models.Persistable[Created], Created any] interface {
	resolvers.AuthDBProvider
	Item() Item
}

type CreateHandlerInput[Item models.Persistable[Created], Created any] struct {
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
	Item models.Persistable[Created],
	Created any,
](ctx context.Context, input *CreateHandlerInput[Item, Created]) (*CreateHandlerOutput[Created], error) {
	return CreateHandlerWithInput(ctx, input)
}

func CreateHandlerWithInput[
	Input CreateInputBody[Item, Created],
	Item models.Persistable[Created],
	Created any,
](ctx context.Context, input Input) (*CreateHandlerOutput[Created], error) {
	created, err := input.Item().Save(input.DB())
	if err = StatusError(err); err != nil {
		return nil, err
	}
	return &CreateHandlerOutput[Created]{Body: created}, nil
}
