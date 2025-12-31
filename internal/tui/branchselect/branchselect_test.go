package branchselect

import (
	"testing"

	"ghlog/internal/models"
)

func TestItemTitle(t *testing.T) {
	tests := []struct {
		name     string
		item     item
		expected string
	}{
		{"regular", item{name: "main", isDefault: false}, "main"},
		{"default", item{name: "main", isDefault: true}, "main (default)"},
		{"featureBranch", item{name: "feature/test", isDefault: false}, "feature/test"},
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
	i := item{name: "main"}
	if got := i.Description(); got != "" {
		t.Errorf("Description() = %q, want empty", got)
	}
}

func TestItemFilterValue(t *testing.T) {
	i := item{name: "develop"}
	if got := i.FilterValue(); got != "develop" {
		t.Errorf("FilterValue() = %q, want %q", got, "develop")
	}
}

func TestModelSubmit(t *testing.T) {
	tests := []struct {
		name           string
		selectedBranch string
		defaultBranch  string
		expected       string
	}{
		{"withSelection", "develop", "main", "develop"},
		{"emptySelection", "", "main", "main"},
		{"featureBranch", "feature/test", "main", "feature/test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				repo:           models.Repository{DefaultBranchName: tt.defaultBranch},
				selectedBranch: tt.selectedBranch,
			}
			msg := m.submit().(DoneMsg)
			if msg.Branch != tt.expected {
				t.Errorf("Branch = %q, want %q", msg.Branch, tt.expected)
			}
		})
	}
}

func TestModelSubmitDefault(t *testing.T) {
	m := Model{}
	msg := m.submitDefault()
	if _, ok := msg.(UseDefaultMsg); !ok {
		t.Errorf("submitDefault() returned %T, want UseDefaultMsg", msg)
	}
}
