package filterform

import (
	"testing"

	"gh-commit-analyzer/internal/models"
)

func TestNew(t *testing.T) {
	m := New()

	if len(m.inputs) != fieldCount {
		t.Errorf("inputs count = %d, want %d", len(m.inputs), fieldCount)
	}
	if m.focused != 0 {
		t.Errorf("focused = %d, want 0", m.focused)
	}
}

func TestFilters(t *testing.T) {
	m := New()
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
	m := New()
	m.inputs[fieldPerPage].SetValue("invalid")

	f := m.Filters()

	if f.PerPage != models.DefaultPerPage {
		t.Errorf("PerPage = %d, want %d", f.PerPage, models.DefaultPerPage)
	}
}

func TestNextField(t *testing.T) {
	m := New()

	m = m.nextField()
	if m.focused != 1 {
		t.Errorf("focused = %d, want 1", m.focused)
	}

	m = m.nextField()
	m = m.nextField()
	m = m.nextField()
	if m.focused != 0 {
		t.Errorf("focused = %d, want 0 (wrap around)", m.focused)
	}
}

func TestPrevField(t *testing.T) {
	m := New()

	m = m.prevField()
	if m.focused != fieldCount-1 {
		t.Errorf("focused = %d, want %d (wrap around)", m.focused, fieldCount-1)
	}

	m = m.prevField()
	if m.focused != fieldCount-2 {
		t.Errorf("focused = %d, want %d", m.focused, fieldCount-2)
	}
}
