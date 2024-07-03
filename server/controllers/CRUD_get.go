package controllers

import (
	"context"
	"darco/proto/models"
	"darco/proto/resolvers"
	"darco/proto/router"

	"github.com/edgedb/edgedb-go"
)

type GetInputInterface[Item any, ID any] interface {
	resolvers.AuthDBProvider
	IdentifierInput[ID] // The identifier of the item to get
}

type GetInput[Item any, ID any] struct {
	Body Item
}

type GetByCodeHandlerInput[Item any] struct {
	resolvers.AuthRequired
	CodeInput
	GetInput[Item, string]
}

type GetByIDHandlerInput[Item any] struct {
	resolvers.AuthRequired
	UUIDInput
	GetInput[Item, edgedb.UUID]
}

type GetHandlerOutput[Item any] struct {
	Body Item
}

func GetHandler[
	OperationInput GetInputInterface[Item, ID],
	Item any,
	ID any,
](
	find models.ItemFinder[ID, Item],
) func(context.Context, OperationInput) (*GetHandlerOutput[Item], error) {
	return func(ctx context.Context, input OperationInput) (*GetHandlerOutput[Item], error) {
		item, err := find(input.DB(), input.Identifier())
		return &GetHandlerOutput[Item]{Body: item}, err
	}
}

func GetByIDHandler[Item any](
	find models.ItemFinder[edgedb.UUID, Item],
) router.Endpoint[
	GetByIDHandlerInput[Item],
	GetHandlerOutput[Item],
] {
	return GetHandler[*GetByIDHandlerInput[Item]](find)
}

func GetByCodeHandler[Item any](
	find models.ItemFinder[string, Item],
) router.Endpoint[
	GetByCodeHandlerInput[Item],
	GetHandlerOutput[Item],
] {
	return GetHandler[*GetByCodeHandlerInput[Item]](find)
}

// Implementation assertions
var _ UpdateInputInterface[models.Updatable[string, any], string, any] = (*UpdateByCodeHandlerInput[models.Updatable[string, any], any])(nil)
var _ UpdateInputInterface[models.Updatable[edgedb.UUID, any], edgedb.UUID, any] = (*UpdateByIDHandlerInput[models.Updatable[edgedb.UUID, any], any])(nil)
