package controllers

import (
	"context"
	"reflect"

	"github.com/lsdch/biome/resolvers"
)

// `FetchItemList` functions retrieve item lists from the database
type FetchItemList[Item any] FetchItem[[]Item]

type ListHandlerInput struct {
	resolvers.AuthResolver
}

type ListHandlerOutput[Item any] struct{ Body []Item }

type ListItemHandler[
	Result any,
	Input resolvers.AuthDBProvider,
	Options any,
] func(ctx context.Context, input Input) (*FetchHandlerOutput[Result], error)

func ListHandler[Input resolvers.AuthDBProvider, Item any](listFn FetchItemList[Item]) ListItemHandler[[]Item, Input, any] {
	return func(ctx context.Context, input Input) (*FetchHandlerOutput[[]Item], error) {
		items, err := listFn(input.DB())
		return handleListItemsResult(items, err)
	}
}

type ListHandlerInputWithOptions[Options any] interface {
	resolvers.AuthDBProvider
	Options() Options
}

func ListHandlerWithOpts[Input ListHandlerInputWithOptions[Options], Result any, Options any](
	listFn FetchItemWithOptions[Result, Options],
) ListItemHandler[Result, Input, Options] {
	return func(ctx context.Context, input Input) (*FetchHandlerOutput[Result], error) {
		items, err := listFn(input.DB(), input.Options())
		return handleListItemsResult(items, err)
	}
}

func handleListItemsResult[Result any](items Result, err error) (*FetchHandlerOutput[Result], error) {
	t := reflect.ValueOf(items)
	if t.Kind() == reflect.Slice && t.Len() == 0 {
		items = reflect.MakeSlice(t.Type().Elem(), 0, 0).Interface().(Result)
	}
	if err = StatusError(err); err != nil {
		return nil, err
	}
	return &FetchHandlerOutput[Result]{Body: items}, nil
}
