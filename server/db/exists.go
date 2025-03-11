package db

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/geldata/gel-go/geltypes"
)

type DBProperty struct {
	Object   string
	Property string
}

func (q DBProperty) Exists(db geltypes.Executor, v string) (geltypes.UUID, bool) {
	var uuid geltypes.UUID
	query := fmt.Sprintf(`select (select %s filter .%s = <str>$0).id`,
		q.Object, q.Property,
	)
	err := db.QuerySingle(context.Background(), query, &uuid, v)
	return uuid, !IsNoData(err)
}

func (q DBProperty) NotExists(db geltypes.Executor, v string) bool {
	var exists bool
	query := fmt.Sprintf(`select exists %s filter .%s = <str>$0`,
		q.Object, q.Property,
	)
	_ = db.QuerySingle(context.Background(), query, &exists)
	return exists
}

func (q DBProperty) ExistAll(db geltypes.Executor, identifiers []string) ([]geltypes.UUID, []InvalidItem) {
	var missings []InvalidItem
	var uuids []geltypes.UUID
	for i, v := range identifiers {
		uuid, ok := q.Exists(db, v)
		if ok {
			uuids = append(uuids, uuid)
		} else {
			missings = append(missings, InvalidItem{Value: v, Position: i})
		}
	}
	if len(missings) > 0 {
		return nil, missings
	}
	return uuids, nil
}

func (q DBProperty) NotExistAll(db geltypes.Executor, identifiers []string) []InvalidItem {
	var invalid []InvalidItem
	for i, v := range identifiers {
		ok := q.NotExists(db, v)
		if !ok {
			invalid = append(invalid, InvalidItem{Value: v, Position: i})
		}
	}
	return invalid
}

type InvalidItem struct {
	Value    string
	Position int
}

func (a *InvalidItem) ErrorDetail(location string) *huma.ErrorDetail {
	return &huma.ErrorDetail{
		Message:  "Item not found",
		Location: fmt.Sprintf("%s[%d]", location, a.Position),
		Value:    a.Value,
	}
}
