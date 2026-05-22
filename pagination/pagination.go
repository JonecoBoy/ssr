package pagination

import (
	"strconv"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
	PageKey      = "page"
	LimitKey     = "limit"
)

type Pagination struct {
	Page  int
	Limit int
}

// Offset returns the number of items to skip based on page and limit.
// Alias for Skip() - use whichever name fits your ORM/database better.
func (p *Pagination) Offset() int {
	return p.Skip()
}

// Skip returns the number of items to skip for pagination.
// Ensures page is at least 1 before calculating.
func (p *Pagination) Skip() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetLimit()
}

// GetLimit returns a validated limit value.
// Returns DefaultLimit if limit is <= 0, MaxLimit if limit exceeds maximum.
func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		return DefaultLimit
	}
	if p.Limit > MaxLimit {
		return MaxLimit
	}
	return p.Limit
}

// GetPage returns a validated page value (minimum 1).
func (p *Pagination) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

// TotalPages calculates total pages based on total items count.
func (p *Pagination) TotalPages(totalItems int64) int {
	limit := p.GetLimit()
	if limit == 0 {
		return 0
	}
	pages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		pages++
	}
	return pages
}

// HasNext returns true if there are more pages after the current one.
func (p *Pagination) HasNext(totalItems int64) bool {
	return p.GetPage() < p.TotalPages(totalItems)
}

// HasPrev returns true if there are pages before the current one.
func (p *Pagination) HasPrev() bool {
	return p.GetPage() > 1
}

// Parse creates a Pagination from a query string map.
// Uses default values if parameters are missing or invalid.
func Parse(query map[string]string) Pagination {
	return ParseWithDefaults(query, DefaultPage, DefaultLimit)
}

// ParseWithDefaults creates a Pagination with custom default values.
func ParseWithDefaults(query map[string]string, defaultPage, defaultLimit int) Pagination {
	pagination := Pagination{
		Page:  defaultPage,
		Limit: defaultLimit,
	}

	if page := query[PageKey]; page != "" {
		if v, err := strconv.Atoi(page); err == nil && v > 0 {
			pagination.Page = v
		}
	}

	if limit := query[LimitKey]; limit != "" {
		if v, err := strconv.Atoi(limit); err == nil && v > 0 {
			pagination.Limit = v
		}
	}

	return pagination
}

// New creates a Pagination with specified page and limit.
func New(page, limit int) Pagination {
	return Pagination{
		Page:  page,
		Limit: limit,
	}
}

// Meta contains pagination metadata for API responses.
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// GetMeta returns pagination metadata for API responses.
func (p *Pagination) GetMeta(totalItems int64) Meta {
	return Meta{
		Page:       p.GetPage(),
		Limit:      p.GetLimit(),
		Total:      totalItems,
		TotalPages: p.TotalPages(totalItems),
		HasNext:    p.HasNext(totalItems),
		HasPrev:    p.HasPrev(),
	}
}
