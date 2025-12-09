package reposelect

import (
	"testing"

	"gh-commit-analyzer/internal/models"
)

func TestItemTitle(t *testing.T) {
	tests := []struct {
		name     string
		item     item
		expected string
	}{
		{
			name:     "unselected",
			item:     item{repo: models.Repository{NameWithOwner: "owner/repo"}, selected: false},
			expected: "[ ] owner/repo",
		},
		{
			name:     "selected",
			item:     item{repo: models.Repository{NameWithOwner: "owner/repo"}, selected: true},
			expected: "[x] owner/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.Title(); got != tt.expected {
				t.Errorf("Title() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestItemDescription(t *testing.T) {
	tests := []struct {
		name     string
		item     item
		contains string
	}{
		{
			name:     "withDescription",
			item:     item{repo: models.Repository{Description: "A test repo"}},
			contains: "A test repo",
		},
		{
			name:     "withoutDescription",
			item:     item{repo: models.Repository{}},
			contains: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.Description()
			if len(got) == 0 {
				t.Error("Description() returned empty string")
			}
		})
	}
}

func TestItemFilterValue(t *testing.T) {
	i := item{repo: models.Repository{NameWithOwner: "owner/repo"}}
	if got := i.FilterValue(); got != "owner/repo" {
		t.Errorf("FilterValue() = %q, want %q", got, "owner/repo")
	}
}

func TestModelSelected(t *testing.T) {
	m := Model{
		selected: map[string]models.Repository{
			"owner/repo1": {NameWithOwner: "owner/repo1"},
			"owner/repo2": {NameWithOwner: "owner/repo2"},
		},
	}

	selected := m.Selected()
	if len(selected) != 2 {
		t.Errorf("Selected() returned %d repos, want 2", len(selected))
	}
}
