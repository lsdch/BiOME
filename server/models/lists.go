package models

type Pagination struct {
	Limit  int `query:"limit" json:"limit,omitzero"`
	Offset int `query:"offset" json:"offset,omitzero"`
}

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

type SortBy[T ~string] struct {
	Key   T         `query:"sort" json:"key,omitempty"`
	Order SortOrder `query:"order" json:"order,omitempty"`
}

type Filter struct {
	SearchTerm string `query:"search" json:"search,omitzero"`
	Owned      bool   `query:"owned" json:"owned"`
}

type PaginatedList[T any] struct {
	Items      []T   `json:"items" gel:"items"`
	TotalCount int64 `json:"total_count" gel:"total_count"`
}
