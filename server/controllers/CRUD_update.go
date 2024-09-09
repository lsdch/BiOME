package controllers

import (
	"context"
	"darco/proto/models"
	"darco/proto/resolvers"
	"darco/proto/router"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type UpdateInputInterface[Item models.Updatable[ID, UpdatedID], ID any, UpdatedID any] interface {
	resolvers.AuthDBProvider
	IdentifierInput[ID] // The identifier of the item to update
	Item() Item         // The update payload
}

type UpdateInput[Item models.Updatable[ID, UpdatedID], ID any, UpdatedID any] struct {
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
	OperationInput UpdateInputInterface[Item, ID, UpdatedID],
	Item models.Updatable[ID, UpdatedID],
	ID any,
	UpdatedID any,
	Updated any,
](
	find models.ItemFinder[UpdatedID, Updated],
) func(context.Context, OperationInput) (*UpdateHandlerOutput[Updated], error) {
	return func(ctx context.Context, input OperationInput) (*UpdateHandlerOutput[Updated], error) {
		uuid, err := input.Item().Update(
			input.DB(),
			input.Identifier(),
		)
		if err != nil {
			logrus.Errorf("Item update failed: %+v", err)
			return nil, huma.Error500InternalServerError("Item update failed", err)
		}

		updated, err := find(input.DB(), uuid)
		if err != nil {
			logrus.Errorf("Failed to retrieve updated item: %v", err)
			return nil, huma.Error500InternalServerError("Failed to retrieve updated item", err)
		}

		logrus.Infof("Item updated: %+v", updated)
		return &UpdateHandlerOutput[Updated]{Body: updated}, nil
	}
}

func UpdateByIDHandler[Item models.Updatable[edgedb.UUID, UpdatedID], Updated any, UpdatedID any](
	find models.ItemFinder[UpdatedID, Updated],
) router.Endpoint[
	UpdateByIDHandlerInput[Item, UpdatedID],
	UpdateHandlerOutput[Updated],
] {
	return UpdateHandler[*UpdateByIDHandlerInput[Item, UpdatedID]](find)
}

func UpdateByCodeHandler[Item models.Updatable[string, UpdatedID], Updated any, UpdatedID any](
	find models.ItemFinder[UpdatedID, Updated],
) router.Endpoint[
	UpdateByCodeHandlerInput[Item, UpdatedID],
	UpdateHandlerOutput[Updated],
] {
	return UpdateHandler[*UpdateByCodeHandlerInput[Item, UpdatedID]](find)
}

// Implementation assertions
var _ UpdateInputInterface[models.Updatable[string, any], string, any] = (*UpdateByCodeHandlerInput[models.Updatable[string, any], any])(nil)
var _ UpdateInputInterface[models.Updatable[edgedb.UUID, any], edgedb.UUID, any] = (*UpdateByIDHandlerInput[models.Updatable[edgedb.UUID, any], any])(nil)
