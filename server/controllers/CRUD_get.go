package controllers

import (
	"context"

	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/edgedb/edgedb-go"
)

type GetInputInterface[Item any, ID any] interface {
	resolvers.AuthDBProvider
	IdentifierInput[ID] // The identifier of the item to get
}

type GetByCodeHandlerInput[Item any] struct {
	resolvers.AuthResolver
	CodeInput
}

type GetByIDHandlerInput[Item any] struct {
	resolvers.AuthResolver
	UUIDInput
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
		if err = StatusError(err); err != nil {
			return nil, err
		}
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
var _ UpdateInputInterface[models.PersistableWithID[string, any], string, any] = (*UpdateByCodeHandlerInput[models.PersistableWithID[string, any], any])(nil)
var _ UpdateInputInterface[models.PersistableWithID[edgedb.UUID, any], edgedb.UUID, any] = (*UpdateByIDHandlerInput[models.PersistableWithID[edgedb.UUID, any], any])(nil)
