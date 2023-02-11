package common

import "math"

type PaginatedResponse[T any] struct {
	List          []T   `json:"list"`
	Page          int   `json:"page"`
	Pages         int   `json:"pages"`
	PageSize      int   `json:"pageSize"`
	TotalElements int64 `json:"totalElements"`
}

type Paginator[T any] struct {
	PageSize   int
	Items      *[]T
	TotalItems int64
}

func (p *Paginator[T]) Get(page int) *PaginatedResponse[T] {
	itemsLenght := float64(p.TotalItems)
	pageSize := float64(p.PageSize)
	pages := int(math.Ceil(itemsLenght / pageSize))

	res := PaginatedResponse[T]{
		List:          *p.Items,
		Page:          page,
		Pages:         pages,
		PageSize:      p.PageSize,
		TotalElements: p.TotalItems,
	}

	return &res
}
