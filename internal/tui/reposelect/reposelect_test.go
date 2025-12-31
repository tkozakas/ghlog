package reposelect

import (
	"testing"

	"ghlog/internal/models"
)

func TestItemTitle(t *testing.T) {
	unselected := make(map[string]models.Repository)
	selected := map[string]models.Repository{
		"owner/repo": {NameWithOwner: "owner/repo"},
	}

	tests := []struct {
		name     string
		item     item
		expected string
	}{
		{
			name:     "unselected",
			item:     item{repo: models.Repository{NameWithOwner: "owner/repo"}, selected: unselected},
			expected: "[ ] owner/repo",
		},
		{
			name:     "selected",
			item:     item{repo: models.Repository{NameWithOwner: "owner/repo"}, selected: selected},
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
	selected := make(map[string]models.Repository)

	tests := []struct {
		name     string
		item     item
		contains string
	}{
		{
			name:     "withDescription",
			item:     item{repo: models.Repository{Description: "A test repo"}, selected: selected},
			contains: "A test repo",
		},
		{
			name:     "withoutDescription",
			item:     item{repo: models.Repository{}, selected: selected},
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
	selected := make(map[string]models.Repository)
	i := item{repo: models.Repository{NameWithOwner: "owner/repo"}, selected: selected}
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

func TestSharedMapBetweenItems(t *testing.T) {
	repos := []models.Repository{
		{NameWithOwner: "foo/bar"},
		{NameWithOwner: "foo/baz"},
	}

	m := New(repos, 80, 24)

	m.selected["foo/baz"] = repos[1]

	items := m.list.Items()
	for _, listItem := range items {
		it := listItem.(item)
		title := it.Title()
		if it.repo.NameWithOwner == "foo/baz" {
			if title[:3] != "[x]" {
				t.Errorf("foo/baz should be checked, got: %s", title)
			}
		} else {
			if title[:3] != "[ ]" {
				t.Errorf("%s should be unchecked, got: %s", it.repo.NameWithOwner, title)
			}
		}
	}
}
