package pagination

import (
	"testing"
)

func TestPagination_Skip(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expected int
	}{
		{"page 1", 1, 20, 0},
		{"page 2", 2, 20, 20},
		{"page 3 with limit 10", 3, 10, 20},
		{"page 0 should default to 1", 0, 20, 0},
		{"negative page should default to 1", -1, 20, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{Page: tt.page, Limit: tt.limit}
			if got := p.Skip(); got != tt.expected {
				t.Errorf("Skip() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPagination_Offset(t *testing.T) {
	p := Pagination{Page: 2, Limit: 10}
	if p.Offset() != p.Skip() {
		t.Error("Offset() should equal Skip()")
	}
}

func TestPagination_GetLimit(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expected int
	}{
		{"normal limit", 20, 20},
		{"zero limit returns default", 0, DefaultLimit},
		{"negative limit returns default", -5, DefaultLimit},
		{"exceeds max returns max", 150, MaxLimit},
		{"at max limit", MaxLimit, MaxLimit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{Page: 1, Limit: tt.limit}
			if got := p.GetLimit(); got != tt.expected {
				t.Errorf("GetLimit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPagination_GetPage(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		expected int
	}{
		{"normal page", 5, 5},
		{"zero page returns 1", 0, 1},
		{"negative page returns 1", -3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{Page: tt.page, Limit: 20}
			if got := p.GetPage(); got != tt.expected {
				t.Errorf("GetPage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPagination_TotalPages(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		totalItems int64
		expected   int
	}{
		{"exact division", 10, 100, 10},
		{"with remainder", 10, 105, 11},
		{"less than one page", 20, 5, 1},
		{"zero items", 20, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{Page: 1, Limit: tt.limit}
			if got := p.TotalPages(tt.totalItems); got != tt.expected {
				t.Errorf("TotalPages() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPagination_HasNext(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		limit      int
		totalItems int64
		expected   bool
	}{
		{"has next page", 1, 10, 50, true},
		{"on last page", 5, 10, 50, false},
		{"beyond last page", 6, 10, 50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{Page: tt.page, Limit: tt.limit}
			if got := p.HasNext(tt.totalItems); got != tt.expected {
				t.Errorf("HasNext() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPagination_HasPrev(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		expected bool
	}{
		{"on first page", 1, false},
		{"on second page", 2, true},
		{"on later page", 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{Page: tt.page, Limit: 20}
			if got := p.HasPrev(); got != tt.expected {
				t.Errorf("HasPrev() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		query         map[string]string
		expectedPage  int
		expectedLimit int
	}{
		{
			"valid page and limit",
			map[string]string{"page": "3", "limit": "50"},
			3, 50,
		},
		{
			"missing page uses default",
			map[string]string{"limit": "30"},
			DefaultPage, 30,
		},
		{
			"missing limit uses default",
			map[string]string{"page": "2"},
			2, DefaultLimit,
		},
		{
			"empty query uses defaults",
			map[string]string{},
			DefaultPage, DefaultLimit,
		},
		{
			"invalid page uses default",
			map[string]string{"page": "abc", "limit": "10"},
			DefaultPage, 10,
		},
		{
			"invalid limit uses default",
			map[string]string{"page": "2", "limit": "xyz"},
			2, DefaultLimit,
		},
		{
			"negative page uses default",
			map[string]string{"page": "-1", "limit": "10"},
			DefaultPage, 10,
		},
		{
			"zero page uses default",
			map[string]string{"page": "0", "limit": "10"},
			DefaultPage, 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parse(tt.query)
			if p.Page != tt.expectedPage {
				t.Errorf("Parse() page = %v, want %v", p.Page, tt.expectedPage)
			}
			if p.Limit != tt.expectedLimit {
				t.Errorf("Parse() limit = %v, want %v", p.Limit, tt.expectedLimit)
			}
		})
	}
}

func TestParseWithDefaults(t *testing.T) {
	query := map[string]string{}
	p := ParseWithDefaults(query, 5, 50)

	if p.Page != 5 {
		t.Errorf("ParseWithDefaults() page = %v, want 5", p.Page)
	}
	if p.Limit != 50 {
		t.Errorf("ParseWithDefaults() limit = %v, want 50", p.Limit)
	}
}

func TestNew(t *testing.T) {
	p := New(3, 25)
	if p.Page != 3 {
		t.Errorf("New() page = %v, want 3", p.Page)
	}
	if p.Limit != 25 {
		t.Errorf("New() limit = %v, want 25", p.Limit)
	}
}

func TestPagination_GetMeta(t *testing.T) {
	p := Pagination{Page: 2, Limit: 10}
	meta := p.GetMeta(55)

	if meta.Page != 2 {
		t.Errorf("GetMeta() Page = %v, want 2", meta.Page)
	}
	if meta.Limit != 10 {
		t.Errorf("GetMeta() Limit = %v, want 10", meta.Limit)
	}
	if meta.Total != 55 {
		t.Errorf("GetMeta() Total = %v, want 55", meta.Total)
	}
	if meta.TotalPages != 6 {
		t.Errorf("GetMeta() TotalPages = %v, want 6", meta.TotalPages)
	}
	if meta.HasNext != true {
		t.Errorf("GetMeta() HasNext = %v, want true", meta.HasNext)
	}
	if meta.HasPrev != true {
		t.Errorf("GetMeta() HasPrev = %v, want true", meta.HasPrev)
	}
}
