package models

import (
	"testing"
	"time"
)

func TestCommitShortSHA(t *testing.T) {
	tests := []struct {
		name     string
		sha      string
		expected string
	}{
		{"fullSha", "abc1234567890", "abc1234"},
		{"exactly7", "abc1234", "abc1234"},
		{"shortSha", "abc", "abc"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Commit{SHA: tt.sha}
			if got := c.ShortSHA(); got != tt.expected {
				t.Errorf("ShortSHA() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestCommitFirstLine(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{"singleLine", "Fix bug", "Fix bug"},
		{"multiLine", "Fix bug\n\nDetailed description", "Fix bug"},
		{"withSpaces", "  Fix bug  \n\nMore", "Fix bug"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Commit{Message: tt.message}
			if got := c.FirstLine(); got != tt.expected {
				t.Errorf("FirstLine() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestCommitHasMultipleLines(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected bool
	}{
		{"singleLine", "Fix bug", false},
		{"multiLine", "Fix bug\n\nDetails", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Commit{Message: tt.message}
			if got := c.HasMultipleLines(); got != tt.expected {
				t.Errorf("HasMultipleLines() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCommitExtraLineCount(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected int
	}{
		{"singleLine", "Fix bug", 0},
		{"twoLines", "Fix bug\nDetails", 1},
		{"fiveLines", "Line1\nLine2\nLine3\nLine4\nLine5", 4},
		{"empty", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Commit{Message: tt.message}
			if got := c.ExtraLineCount(); got != tt.expected {
				t.Errorf("ExtraLineCount() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestCommitFormattedDate(t *testing.T) {
	c := Commit{Date: time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)}
	expected := "2024-06-15 14:30"
	if got := c.FormattedDate(); got != expected {
		t.Errorf("FormattedDate() = %q, want %q", got, expected)
	}
}

func TestCommitFormattedDateZero(t *testing.T) {
	c := Commit{}
	if got := c.FormattedDate(); got != "unknown" {
		t.Errorf("FormattedDate() = %q, want %q", got, "unknown")
	}
}

func TestCommitAuthorWithEmail(t *testing.T) {
	tests := []struct {
		name     string
		author   string
		email    string
		expected string
	}{
		{"withEmail", "John", "john@example.com", "John <john@example.com>"},
		{"withoutEmail", "John", "", "John"},
		{"empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Commit{Author: tt.author, Email: tt.email}
			if got := c.AuthorWithEmail(); got != tt.expected {
				t.Errorf("AuthorWithEmail() = %q, want %q", got, tt.expected)
			}
		})
	}
}
