package controllers

import (
	"context"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/resolvers"
)

// `FetchItem` functions retrieve a single item from the database
type FetchItem[Item any] func(db geltypes.Executor) (Item, error)

type FetchHandlerOutput[Item any] struct{ Body Item }

type FetchItemHandler[
	Item any,
	Input resolvers.AuthDBProvider,
] func(ctx context.Context, input Input) (*FetchHandlerOutput[Item], error)

func FetchHandler[Input resolvers.AuthDBProvider, Item any](fetchFn FetchItem[Item]) FetchItemHandler[Item, Input] {
	return func(ctx context.Context, input Input) (*FetchHandlerOutput[Item], error) {
		item, err := fetchFn(input.DB())
		if err = StatusError(err); err != nil {
			return nil, err
		}
		return &FetchHandlerOutput[Item]{Body: item}, nil
	}
}
