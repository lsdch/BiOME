package controllers

import (
	"context"
	"darco/proto/resolvers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

type FetchItemList[Item any] func(db edgedb.Executor) ([]Item, error)

type ListHandlerInput struct {
	resolvers.AuthResolver
}

type ListHandlerOutput[Item any] struct{ Body []Item }

type ListItemHandler[Item any] func(ctx context.Context, input *ListHandlerInput) (*ListHandlerOutput[Item], error)

func ListHandler[Item any](listFn FetchItemList[Item]) ListItemHandler[Item] {
	return func(ctx context.Context, input *ListHandlerInput) (*ListHandlerOutput[Item], error) {
		items, err := listFn(input.DB())
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to retrieve item list", err)
		}
		return &ListHandlerOutput[Item]{Body: items}, nil
	}
}
