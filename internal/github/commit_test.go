package github

import (
	"testing"

	"github.com/tkozakas/gh-log/internal/models"
)

func TestBuildCommitsEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		owner    string
		repo     string
		branch   string
		filters  models.FilterOptions
		page     int
		expected string
	}{
		{
			name:     "basic",
			owner:    "owner",
			repo:     "repo",
			branch:   "",
			filters:  models.FilterOptions{PerPage: 50},
			page:     1,
			expected: "repos/owner/repo/commits?per_page=50&page=1",
		},
		{
			name:     "withBranch",
			owner:    "owner",
			repo:     "repo",
			branch:   "main",
			filters:  models.FilterOptions{PerPage: 50},
			page:     1,
			expected: "repos/owner/repo/commits?per_page=50&page=1&sha=main",
		},
		{
			name:     "withDateFrom",
			owner:    "owner",
			repo:     "repo",
			branch:   "",
			filters:  models.FilterOptions{PerPage: 50, DateFrom: "2024-01-01"},
			page:     1,
			expected: "repos/owner/repo/commits?per_page=50&page=1&since=2024-01-01T00:00:00Z",
		},
		{
			name:     "withDateTo",
			owner:    "owner",
			repo:     "repo",
			branch:   "",
			filters:  models.FilterOptions{PerPage: 50, DateTo: "2024-12-31"},
			page:     1,
			expected: "repos/owner/repo/commits?per_page=50&page=1&until=2024-12-31T23:59:59Z",
		},
		{
			name:     "withAuthor",
			owner:    "owner",
			repo:     "repo",
			branch:   "",
			filters:  models.FilterOptions{PerPage: 50, Author: "john"},
			page:     1,
			expected: "repos/owner/repo/commits?per_page=50&page=1&author=john",
		},
		{
			name:   "allFilters",
			owner:  "owner",
			repo:   "repo",
			branch: "develop",
			filters: models.FilterOptions{
				PerPage:  100,
				DateFrom: "2024-01-01",
				DateTo:   "2024-06-30",
				Author:   "jane",
			},
			page:     2,
			expected: "repos/owner/repo/commits?per_page=100&page=2&sha=develop&since=2024-01-01T00:00:00Z&until=2024-06-30T23:59:59Z&author=jane",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCommitsEndpoint(tt.owner, tt.repo, tt.branch, tt.filters, tt.page)
			if got != tt.expected {
				t.Errorf("buildCommitsEndpoint() = %q, want %q", got, tt.expected)
			}
		})
	}
}
