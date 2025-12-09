package models

import "testing"

func TestNewFilterOptions(t *testing.T) {
	f := NewFilterOptions()
	if f.PerPage != DefaultPerPage {
		t.Errorf("PerPage = %d, want %d", f.PerPage, DefaultPerPage)
	}
}

func TestFilterOptionsValidate(t *testing.T) {
	tests := []struct {
		name     string
		perPage  int
		expected int
	}{
		{"zeroSetsDefault", 0, DefaultPerPage},
		{"negativeSetsDefault", -1, DefaultPerPage},
		{"overMaxSetsMax", 200, MaxPerPage},
		{"validUnchanged", 75, 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FilterOptions{PerPage: tt.perPage}
			f.Validate()
			if f.PerPage != tt.expected {
				t.Errorf("PerPage = %d, want %d", f.PerPage, tt.expected)
			}
		})
	}
}

func TestFilterOptionsHasAnyFilter(t *testing.T) {
	tests := []struct {
		name     string
		filter   FilterOptions
		expected bool
	}{
		{"noFilters", FilterOptions{}, false},
		{"dateFrom", FilterOptions{DateFrom: "2024-01-01"}, true},
		{"dateTo", FilterOptions{DateTo: "2024-12-31"}, true},
		{"author", FilterOptions{Author: "john"}, true},
		{"allFilters", FilterOptions{DateFrom: "2024-01-01", DateTo: "2024-12-31", Author: "john"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.HasAnyFilter(); got != tt.expected {
				t.Errorf("HasAnyFilter() = %v, want %v", got, tt.expected)
			}
		})
	}
}
