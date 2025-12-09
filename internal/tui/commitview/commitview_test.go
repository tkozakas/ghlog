package commitview

import (
	"testing"

	"gh-commit-analyzer/internal/models"
)

func TestModelCountCommits(t *testing.T) {
	tests := []struct {
		name        string
		repoCommits []models.RepoCommits
		expected    int
	}{
		{
			name:        "empty",
			repoCommits: []models.RepoCommits{},
			expected:    0,
		},
		{
			name: "singleRepo",
			repoCommits: []models.RepoCommits{
				{Commits: make([]models.Commit, 5)},
			},
			expected: 5,
		},
		{
			name: "multipleRepos",
			repoCommits: []models.RepoCommits{
				{Commits: make([]models.Commit, 3)},
				{Commits: make([]models.Commit, 7)},
				{Commits: make([]models.Commit, 2)},
			},
			expected: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{repoCommits: tt.repoCommits}
			if got := m.countCommits(); got != tt.expected {
				t.Errorf("countCommits() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestModelMoveCursor(t *testing.T) {
	tests := []struct {
		name         string
		totalCommits int
		startCursor  int
		delta        int
		expected     int
	}{
		{"moveDown", 10, 0, 1, 1},
		{"moveUp", 10, 5, -1, 4},
		{"clampAtZero", 10, 0, -1, 0},
		{"clampAtMax", 10, 9, 1, 9},
		{"moveMultiple", 10, 3, 3, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				totalCommits: tt.totalCommits,
				cursor:       tt.startCursor,
			}
			m.moveCursor(tt.delta)
			if m.cursor != tt.expected {
				t.Errorf("cursor = %d, want %d", m.cursor, tt.expected)
			}
		})
	}
}

func TestModelToggleExpanded(t *testing.T) {
	m := Model{
		expanded: make(map[int]bool),
		cursor:   2,
	}

	m.toggleExpanded()
	if !m.expanded[2] {
		t.Error("expected commit 2 to be expanded")
	}

	m.toggleExpanded()
	if m.expanded[2] {
		t.Error("expected commit 2 to be collapsed")
	}
}
