package filterform

import (
	"testing"

	"gh-commit-analyzer/internal/models"
)

func TestNew(t *testing.T) {
	m := New(nil)

	if len(m.inputs) != fieldCountBase {
		t.Errorf("inputs count = %d, want %d", len(m.inputs), fieldCountBase)
	}
	if m.focused != 0 {
		t.Errorf("focused = %d, want 0", m.focused)
	}
}

func TestNewWithBranches(t *testing.T) {
	repos := []RepoBranches{
		{Repo: models.Repository{NameWithOwner: "org/repo1", DefaultBranchName: "main"}, Branches: []string{"main", "dev"}},
		{Repo: models.Repository{NameWithOwner: "org/repo2", DefaultBranchName: "master"}, Branches: []string{"master"}},
	}
	m := New(repos)

	if m.fieldCount != fieldCountBase+2 {
		t.Errorf("fieldCount = %d, want %d", m.fieldCount, fieldCountBase+2)
	}
	if len(m.branchIdx) != 2 {
		t.Errorf("branchIdx length = %d, want 2", len(m.branchIdx))
	}
}

func TestFilters(t *testing.T) {
	m := New(nil)
	m.inputs[fieldDateFrom].SetValue("2024-01-01")
	m.inputs[fieldDateTo].SetValue("2024-06-30")
	m.inputs[fieldAuthor].SetValue("john")
	m.inputs[fieldPerPage].SetValue("25")

	f := m.Filters()

	if f.DateFrom != "2024-01-01" {
		t.Errorf("DateFrom = %q, want %q", f.DateFrom, "2024-01-01")
	}
	if f.DateTo != "2024-06-30" {
		t.Errorf("DateTo = %q, want %q", f.DateTo, "2024-06-30")
	}
	if f.Author != "john" {
		t.Errorf("Author = %q, want %q", f.Author, "john")
	}
	if f.PerPage != 25 {
		t.Errorf("PerPage = %d, want %d", f.PerPage, 25)
	}
}

func TestFiltersValidation(t *testing.T) {
	m := New(nil)
	m.inputs[fieldPerPage].SetValue("invalid")

	f := m.Filters()

	if f.PerPage != models.DefaultPerPage {
		t.Errorf("PerPage = %d, want %d", f.PerPage, models.DefaultPerPage)
	}
}

func TestBranches(t *testing.T) {
	repos := []RepoBranches{
		{Repo: models.Repository{NameWithOwner: "org/repo1", DefaultBranchName: "main"}, Branches: []string{"main", "dev"}},
	}
	m := New(repos)

	branches := m.Branches()

	if branches["org/repo1"] != "main" {
		t.Errorf("branch = %q, want %q", branches["org/repo1"], "main")
	}

	m.branchIdx[0] = 1
	branches = m.Branches()

	if branches["org/repo1"] != "dev" {
		t.Errorf("branch = %q, want %q", branches["org/repo1"], "dev")
	}
}

func TestNextField(t *testing.T) {
	m := New(nil)

	m = m.nextField()
	if m.focused != 1 {
		t.Errorf("focused = %d, want 1", m.focused)
	}

	// fieldCountBase is 5 (dateFrom, dateTo, author, perPage, semanticQuery)
	m = m.nextField()
	m = m.nextField()
	m = m.nextField()
	m = m.nextField()
	if m.focused != 0 {
		t.Errorf("focused = %d, want 0 (wrap around)", m.focused)
	}
}

func TestPrevField(t *testing.T) {
	m := New(nil)

	m = m.prevField()
	if m.focused != m.fieldCount-1 {
		t.Errorf("focused = %d, want %d (wrap around)", m.focused, m.fieldCount-1)
	}

	m = m.prevField()
	if m.focused != m.fieldCount-2 {
		t.Errorf("focused = %d, want %d", m.focused, m.fieldCount-2)
	}
}
