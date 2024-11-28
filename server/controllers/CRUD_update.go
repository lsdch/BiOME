package controllers

import (
	"context"
	"darco/proto/models"
	"darco/proto/resolvers"

	"github.com/edgedb/edgedb-go"
)

type UpdateInputInterface[Item models.Updatable[ID, Updated], ID any, Updated any] interface {
	resolvers.AuthDBProvider
	IdentifierInput[ID] // The identifier of the item to update
	Item() Item         // The update payload
}

type UpdateInput[Item models.Updatable[ID, Updated], ID any, Updated any] struct {
	Body Item
}

func (i UpdateInput[Item, ID, UpdatedID]) Item() Item {
	return i.Body
}

type UpdateByCodeHandlerInput[Item models.Updatable[string, UpdatedID], UpdatedID any] struct {
	resolvers.AuthRequired
	CodeInput
	UpdateInput[Item, string, UpdatedID]
}

type UpdateByIDHandlerInput[Item models.Updatable[edgedb.UUID, UpdatedID], UpdatedID any] struct {
	resolvers.AuthRequired
	UUIDInput
	UpdateInput[Item, edgedb.UUID, UpdatedID]
}

type UpdateHandlerOutput[Updated any] struct {
	Body Updated
}

func UpdateHandler[
	OperationInput UpdateInputInterface[Item, ID, Updated],
	Item models.Updatable[ID, Updated],
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

func UpdateByIDHandler[Item models.Updatable[edgedb.UUID, Updated], Updated any](ctx context.Context, input *UpdateByIDHandlerInput[Item, Updated]) (*UpdateHandlerOutput[Updated], error) {
	return UpdateHandler(ctx, input)
}

func UpdateByCodeHandler[Item models.Updatable[string, Updated], Updated any](ctx context.Context, input *UpdateByCodeHandlerInput[Item, Updated]) (*UpdateHandlerOutput[Updated], error) {
	return UpdateHandler(ctx, input)
}

// Implementation assertions
var _ UpdateInputInterface[models.Updatable[string, any], string, any] = (*UpdateByCodeHandlerInput[models.Updatable[string, any], any])(nil)
var _ UpdateInputInterface[models.Updatable[edgedb.UUID, any], edgedb.UUID, any] = (*UpdateByIDHandlerInput[models.Updatable[edgedb.UUID, any], any])(nil)
