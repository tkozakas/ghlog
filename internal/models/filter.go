package models
// nigger

const (
	DefaultPerPage = 50
	MaxPerPage     = 100
)

type FilterOptions struct {
	DateFrom      string
	DateTo        string
	Author        string
	PerPage       int
	SemanticQuery string
}

type BranchSelection struct {
	RepoName string
	Branch   string
}

func NewFilterOptions() FilterOptions {
	return FilterOptions{PerPage: DefaultPerPage}
}

func (f *FilterOptions) Validate() {
	if f.PerPage <= 0 {
		f.PerPage = DefaultPerPage
	}
	if f.PerPage > MaxPerPage {
		f.PerPage = MaxPerPage
	}
}

func (f FilterOptions) HasAnyFilter() bool {
	return f.hasDateFilter() || f.hasAuthorFilter() || f.HasSemanticFilter()
}

func (f FilterOptions) HasSemanticFilter() bool {
	return f.SemanticQuery != ""
}

func (f FilterOptions) hasDateFilter() bool {
	return f.DateFrom != "" || f.DateTo != ""
}

func (f FilterOptions) hasAuthorFilter() bool {
	return f.Author != ""
}
