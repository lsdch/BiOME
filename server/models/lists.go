package models

type Pagination struct {
	Limit  int `query:"limit" json:"limit,omitzero"`
	Offset int `query:"offset" json:"offset,omitzero"`
}

type Sorting struct {
	SortBy []string `query:"sort_by" json:"sort_by,omitempty"`
}

type Filter struct {
	SearchTerm string `query:"search" json:"search,omitzero"`
	Owned      bool   `query:"owned" json:"owned"`
}

type PaginatedList[T any] struct {
	Items      []T   `json:"items" gel:"items"`
	TotalCount int64 `json:"total_count" gel:"total_count"`
}
