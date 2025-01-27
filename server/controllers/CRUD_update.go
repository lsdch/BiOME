package controllers

import (
	"context"

	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/resolvers"

	"github.com/edgedb/edgedb-go"
)

type UpdateInputInterface[Item models.PersistableWithID[ID, Updated], ID any, Updated any] interface {
	resolvers.AuthDBProvider
	IdentifierInput[ID] // The identifier of the item to update
	Item() Item         // The update payload
}

type UpdateInput[Item models.PersistableWithID[ID, Updated], ID any, Updated any] struct {
	Body Item
}

func (i UpdateInput[Item, ID, UpdatedID]) Item() Item {
	return i.Body
}

type UpdateByCodeHandlerInput[Item models.PersistableWithID[string, UpdatedID], UpdatedID any] struct {
	resolvers.AuthRequired
	CodeInput
	UpdateInput[Item, string, UpdatedID]
}

type UpdateByIDHandlerInput[Item models.PersistableWithID[edgedb.UUID, UpdatedID], UpdatedID any] struct {
	resolvers.AuthRequired
	UUIDInput
	UpdateInput[Item, edgedb.UUID, UpdatedID]
}

type UpdateHandlerOutput[Updated any] struct {
	Body Updated
}

func UpdateHandler[
	OperationInput UpdateInputInterface[Item, ID, Updated],
	Item models.PersistableWithID[ID, Updated],
	ID any,
	Updated any,
](ctx context.Context, input OperationInput) (*UpdateHandlerOutput[Updated], error) {
	updated, err := input.Item().Save(
		input.DB(),
		input.Identifier(),
	)
	if err = StatusError(err); err != nil {
		return nil, err
	}

	return &UpdateHandlerOutput[Updated]{Body: updated}, nil
}

func UpdateByIDHandler[Item models.PersistableWithID[edgedb.UUID, Updated], Updated any](ctx context.Context, input *UpdateByIDHandlerInput[Item, Updated]) (*UpdateHandlerOutput[Updated], error) {
	return UpdateHandler(ctx, input)
}

func UpdateByCodeHandler[Item models.PersistableWithID[string, Updated], Updated any](ctx context.Context, input *UpdateByCodeHandlerInput[Item, Updated]) (*UpdateHandlerOutput[Updated], error) {
	return UpdateHandler(ctx, input)
}

// Implementation assertions
var _ UpdateInputInterface[models.PersistableWithID[string, any], string, any] = (*UpdateByCodeHandlerInput[models.PersistableWithID[string, any], any])(nil)
var _ UpdateInputInterface[models.PersistableWithID[edgedb.UUID, any], edgedb.UUID, any] = (*UpdateByIDHandlerInput[models.PersistableWithID[edgedb.UUID, any], any])(nil)
