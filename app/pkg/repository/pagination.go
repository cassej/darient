package repository

import (
	"fmt"
	"math"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int // 1-based page number
	PageSize int // number of items per page
}

// PaginatedResult holds paginated data
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// DefaultPagination returns default pagination params
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:     1,
		PageSize: 20,
	}
}

// NewPaginationParams creates validated pagination params
func NewPaginationParams(page, pageSize int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset calculates the SQL OFFSET value
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the SQL LIMIT value
func (p PaginationParams) Limit() int {
	return p.PageSize
}

// ToSQL returns LIMIT and OFFSET as SQL string
func (p PaginationParams) ToSQL() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit(), p.Offset())
}

// NewPaginatedResult creates a paginated result
func NewPaginatedResult[T any](items []T, total int64, params PaginationParams) PaginatedResult[T] {
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	
	return PaginatedResult[T]{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}
}

// HasNextPage returns true if there are more pages
func (r PaginatedResult[T]) HasNextPage() bool {
	return r.Page < r.TotalPages
}

// HasPrevPage returns true if there are previous pages
func (r PaginatedResult[T]) HasPrevPage() bool {
	return r.Page > 1
}
