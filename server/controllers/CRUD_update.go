package controllers

import (
	"context"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/resolvers"
)

type UpdateInputInterface[Item models.PersistableWithID[ID, Updated], ID any, Updated any] interface {
	resolvers.AuthDBProvider
	IdentifierInput[ID] // The identifier of the item to update
	Item() Item         // The update payload
}

type UpdateInput[Item models.PersistableWithID[ID, Updated], ID any, Updated any] struct {
	Body Item
}

func (i UpdateInput[Item, ID, Updated]) Item() Item {
	return i.Body
}

type UpdateByCodeHandlerInput[Item models.PersistableWithID[string, Updated], Updated any] struct {
	resolvers.AuthRequired
	CodeInput
	UpdateInput[Item, string, Updated]
}

type UpdateByIDHandlerInput[Item models.PersistableWithID[geltypes.UUID, Updated], Updated any] struct {
	resolvers.AuthRequired
	UUIDInput
	UpdateInput[Item, geltypes.UUID, Updated]
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

func UpdateByIDHandler[Item models.PersistableWithID[geltypes.UUID, Updated], Updated any](ctx context.Context, input *UpdateByIDHandlerInput[Item, Updated]) (*UpdateHandlerOutput[Updated], error) {
	return UpdateHandler(ctx, input)
}

func UpdateByCodeHandler[Item models.PersistableWithID[string, Updated], Updated any](ctx context.Context, input *UpdateByCodeHandlerInput[Item, Updated]) (*UpdateHandlerOutput[Updated], error) {
	return UpdateHandler(ctx, input)
}

// Implementation assertions
var _ UpdateInputInterface[models.PersistableWithID[string, any], string, any] = (*UpdateByCodeHandlerInput[models.PersistableWithID[string, any], any])(nil)
var _ UpdateInputInterface[models.PersistableWithID[geltypes.UUID, any], geltypes.UUID, any] = (*UpdateByIDHandlerInput[models.PersistableWithID[geltypes.UUID, any], any])(nil)
