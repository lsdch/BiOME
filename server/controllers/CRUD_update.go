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

type UpdateInputInterface[Item models.Updatable[ID, Updated], ID any, Updated any] interface {
	resolvers.AuthDBProvider
	IdentifierInput[ID] // The identifier of the item to update
	Item() Item         // The update payload
}

type UpdateInput[Item models.Updatable[ID, Updated], ID any, Updated any] struct {
	Body Item
}

func (i UpdateInput[Item, ID, Updated]) Item() Item {
	return i.Body
}

type UpdateByCodeHandlerInput[Item models.Updatable[string, Updated], Updated any] struct {
	resolvers.AuthRequired
	CodeInput
	UpdateInput[Item, string, Updated]
}

type UpdateByIDHandlerInput[Item models.Updatable[edgedb.UUID, Updated], Updated any] struct {
	resolvers.AuthRequired
	UUIDInput
	UpdateInput[Item, edgedb.UUID, Updated]
}

type UpdateHandlerOutput[Updated any] struct {
	Body Updated
}

type UpdateEndpoint[
	Input UpdateInputInterface[Item, ID, Updated],
	Item models.Updatable[ID, Updated],
	ID any,
	Updated any,
] func(context.Context, Input) (*UpdateHandlerOutput[Updated], error)

func UpdateHandler[
	Item models.Updatable[ID, Updated],
	OperationInput UpdateInputInterface[Item, ID, Updated],
	ID any,
	Updated any,
](
	find models.ItemFinder[edgedb.UUID, Updated],
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

func UpdateByIDHandler[Item models.Updatable[edgedb.UUID, Updated], Updated any](
	find models.ItemFinder[edgedb.UUID, Updated],
) router.Endpoint[
	UpdateByIDHandlerInput[Item, Updated],
	UpdateHandlerOutput[Updated],
] {
	return UpdateHandler[Item, *UpdateByIDHandlerInput[Item, Updated]](find)
}

func UpdateByCodeHandler[Item models.Updatable[string, Updated], Updated any](
	find models.ItemFinder[edgedb.UUID, Updated],
) router.Endpoint[
	UpdateByCodeHandlerInput[Item, Updated],
	UpdateHandlerOutput[Updated],
] {
	return UpdateHandler[Item, *UpdateByCodeHandlerInput[Item, Updated]](find)
}

// Implementation assertions
var _ UpdateInputInterface[models.Updatable[string, any], string, any] = (*UpdateByCodeHandlerInput[models.Updatable[string, any], any])(nil)
var _ UpdateInputInterface[models.Updatable[edgedb.UUID, any], edgedb.UUID, any] = (*UpdateByIDHandlerInput[models.Updatable[edgedb.UUID, any], any])(nil)
