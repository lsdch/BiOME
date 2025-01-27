package controllers

import (
	"context"

	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type ItemDelete[ID any, Item any] func(db edgedb.Executor, id ID) (Item, error)

type DeleteHandlerInput[ID any] interface {
	resolvers.AuthDBProvider
	IdentifierInput[ID]
}

type DeleteHandlerOutput[Deleted any] struct {
	Body Deleted
}

func DeleteHandler[
	OperationInput DeleteHandlerInput[ID],
	Item any,
	ID any,
](deleteItem ItemDelete[ID, Item]) func(context.Context, OperationInput) (*DeleteHandlerOutput[Item], error) {
	return func(ctx context.Context, input OperationInput) (*DeleteHandlerOutput[Item], error) {
		deleted, err := deleteItem(input.DB(), input.Identifier())
		if err = StatusError(err); err != nil {
			return nil, err
		}
		logrus.Debugf("Deleted item %+v", deleted)
		return &DeleteHandlerOutput[Item]{
			Body: deleted,
		}, nil
	}
}

type DeleteByCodeHandlerInput struct {
	resolvers.AuthRequired
	CodeInput
}

func DeleteByCodeHandler[Item any](
	deleteItem ItemDelete[string, Item],
) router.Endpoint[DeleteByCodeHandlerInput, DeleteHandlerOutput[Item]] {
	return DeleteHandler[*DeleteByCodeHandlerInput](deleteItem)
}

type DeleteByIDHandlerInput struct {
	resolvers.AuthRequired
	UUIDInput
}

func DeleteByIDHandler[Item any](
	deleteItem ItemDelete[edgedb.UUID, Item],
) router.Endpoint[DeleteByIDHandlerInput, DeleteHandlerOutput[Item]] {
	return DeleteHandler[*DeleteByIDHandlerInput](deleteItem)
}
