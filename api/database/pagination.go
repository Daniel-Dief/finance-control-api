package database

// Pagination defines the page-based pagination parameters used by the list
// repositories. Page is 1-based and Limit is the maximum number of items per
// page.
type Pagination struct {
	Page  int
	Limit int
}

// Offset returns the number of rows to skip for the current page.
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

// PaginatedResult wraps a page of results together with the pagination
// metadata (total count and total number of pages).
type PaginatedResult[T any] struct {
	Data       []T `json:"data"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

// newPaginatedResult builds a PaginatedResult from the given page data, the
// pagination settings and the total number of matching rows.
func newPaginatedResult[T any](data []T, pag Pagination, total int) PaginatedResult[T] {
	totalPages := 0
	if pag.Limit > 0 {
		totalPages = (total + pag.Limit - 1) / pag.Limit
	}
	return PaginatedResult[T]{
		Data:       data,
		Total:      total,
		Page:       pag.Page,
		Limit:      pag.Limit,
		TotalPages: totalPages,
	}
}
