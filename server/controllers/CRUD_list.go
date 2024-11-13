package controllers

import (
	"context"
	"darco/proto/resolvers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

// `FetchItemList` functions retrieve item lists from the database
type FetchItemList[Item any] func(db edgedb.Executor) ([]Item, error)

type ListHandlerInput struct {
	resolvers.AuthResolver
}

type ListHandlerOutput[Item any] struct{ Body []Item }

type ListItemHandler[Item any, Input resolvers.AuthDBProvider] func(ctx context.Context, input Input) (*ListHandlerOutput[Item], error)

func ListHandler[Input resolvers.AuthDBProvider, Item any](listFn FetchItemList[Item]) ListItemHandler[Item, Input] {
	return func(ctx context.Context, input Input) (*ListHandlerOutput[Item], error) {
		items, err := listFn(input.DB())
		if len(items) == 0 {
			items = []Item{}
		}
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to retrieve item list", err)
		}
		return &ListHandlerOutput[Item]{Body: items}, nil
	}
}
