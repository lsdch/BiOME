package models

import (
	. "github.com/go-jet/jet/v2/postgres"
)

type Pagination struct {
	Limit  int32 `query:"limit" json:"limit,omitzero"`
	Offset int32 `query:"offset" json:"offset,omitzero"`
}

type SortOrder string

//generate:enum
const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

type SortKey interface {
	Column() Column
}

type SortBy[T SortKey] struct {
	Key   Optional[T] `query:"sort" json:"key,omitempty"`
	Order SortOrder   `query:"sort_direction" json:"order,omitempty"`
}

func (s SortBy[T]) IsSet() bool {
	return s.Key.IsSet
}

func (s SortBy[T]) ToOrderByClause() OrderByClause {
	column := (s.Key).Value.Column()
	switch s.Order {
	case SortDesc:
		return column.DESC().NULLS_LAST()
	default:
		return column.ASC().NULLS_LAST()
	}
}

type Filter struct {
	SearchTerm string `query:"search" json:"search,omitzero"`
	Owned      bool   `query:"owned" json:"owned"`
}

type PaginatedList[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
}
