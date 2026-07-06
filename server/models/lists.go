package models

import (
	. "github.com/go-jet/jet/v2/postgres"
)

type Pagination struct {
	Limit  int32 `query:"limit" json:"limit,omitzero"`
	Offset int32 `query:"offset" json:"offset,omitzero"`
}

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

type SortKey interface {
	Column() Column
}

type SortBy[T SortKey] struct {
	Key      T    `query:"sort" json:"key,omitempty"`
	OrderAsc bool `query:"sortAsc" json:"order,omitempty" doc:"ASC if true, default is DESC"`
}

func (s SortBy[T]) ToOrderByClause() OrderByClause {
	if s.OrderAsc == false {
		return s.Key.Column().DESC()
	}
	return s.Key.Column().ASC()
}

type Filter struct {
	SearchTerm string `query:"search" json:"search,omitzero"`
	Owned      bool   `query:"owned" json:"owned"`
}

type PaginatedList[T any] struct {
	Items      []T   `json:"items" gel:"items"`
	TotalCount int64 `json:"total_count" gel:"total_count"`
}
