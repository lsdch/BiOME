package db

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

type DBProperty struct {
	Object   string
	Property string
}

func (q DBProperty) Exists(db edgedb.Executor, v string) (edgedb.UUID, bool) {
	var uuid edgedb.UUID
	query := fmt.Sprintf(`select %s.id filter .%s = <str>$0`,
		q.Object, q.Property,
	)
	err := db.QuerySingle(context.Background(), query, &uuid)
	return uuid, !IsNoData(err)
}

type AbsentItem struct {
	Value    string
	Position int
}

func (a *AbsentItem) ErrorDetail(location string) *huma.ErrorDetail {
	return &huma.ErrorDetail{
		Message:  "Item not found",
		Location: fmt.Sprintf("%s[%d]", location, a.Position),
		Value:    a.Value,
	}
}

func (q DBProperty) ExistAll(db edgedb.Executor, identifiers []string) ([]edgedb.UUID, []AbsentItem) {
	var missings []AbsentItem
	var uuids []edgedb.UUID
	for i, v := range identifiers {
		uuid, ok := q.Exists(db, v)
		if ok {
			uuids = append(uuids, uuid)
		} else {
			missings = append(missings, AbsentItem{Value: v, Position: i})
		}
	}
	if len(missings) > 0 {
		return nil, missings
	}
	return uuids, nil
}
