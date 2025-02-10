package controllers

import (
	"context"

	"github.com/lsdch/biome/resolvers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

// `FetchItemList` functions retrieve item lists from the database
type FetchItemList[Item any] func(db edgedb.Executor) ([]Item, error)

type FetchItemListWithOptions[Item any, Options any] func(db edgedb.Executor, options Options) ([]Item, error)

type ListHandlerInput struct {
	resolvers.AuthResolver
}

type ListHandlerInputWithOptions[Options any] interface {
	resolvers.AuthDBProvider
	Options() Options
}

type ListHandlerOutput[Item any] struct{ Body []Item }

type ListItemHandler[
	Item any,
	Input resolvers.AuthDBProvider,
	Options any,
] func(ctx context.Context, input Input) (*ListHandlerOutput[Item], error)

func ListHandlerWithOpts[Input ListHandlerInputWithOptions[Options], Item any, Options any](
	listFn FetchItemListWithOptions[Item, Options],
) ListItemHandler[Item, Input, Options] {
	return func(ctx context.Context, input Input) (*ListHandlerOutput[Item], error) {
		items, err := listFn(input.DB(), input.Options())
		if len(items) == 0 {
			items = []Item{}
		}
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to retrieve item list", err)
		}
		return &ListHandlerOutput[Item]{Body: items}, nil
	}
}

func ListHandler[Input resolvers.AuthDBProvider, Item any](listFn FetchItemList[Item]) ListItemHandler[Item, Input, any] {
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
