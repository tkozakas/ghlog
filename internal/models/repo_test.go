package models

import (
	"testing"
	"time"
)

func TestRepositoryOwner(t *testing.T) {
	tests := []struct {
		name          string
		nameWithOwner string
		expected      string
	}{
		{"standard", "owner/repo", "owner"},
		{"withDots", "my.org/my.repo", "my.org"},
		{"empty", "", ""},
		{"noSlash", "repo", "repo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{NameWithOwner: tt.nameWithOwner}
			if got := r.Owner(); got != tt.expected {
				t.Errorf("Owner() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRepositoryRepoName(t *testing.T) {
	tests := []struct {
		name          string
		nameWithOwner string
		fallbackName  string
		expected      string
	}{
		{"standard", "owner/repo", "", "repo"},
		{"withDots", "my.org/my.repo", "", "my.repo"},
		{"noSlashUsesFallback", "", "fallback", "fallback"},
		{"noSlashNoFallback", "justname", "justname", "justname"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{NameWithOwner: tt.nameWithOwner, Name: tt.fallbackName}
			if got := r.RepoName(); got != tt.expected {
				t.Errorf("RepoName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRepositoryTimeSincePush(t *testing.T) {
	tests := []struct {
		name     string
		pushedAt time.Time
		expected string
	}{
		{"zeroTime", time.Time{}, "unknown"},
		{"justNow", time.Now().Add(-30 * time.Second), "just now"},
		{"minutes", time.Now().Add(-5 * time.Minute), "5 minutes ago"},
		{"oneMinute", time.Now().Add(-1 * time.Minute), "1 minute ago"},
		{"hours", time.Now().Add(-3 * time.Hour), "3 hours ago"},
		{"oneHour", time.Now().Add(-1 * time.Hour), "1 hour ago"},
		{"days", time.Now().Add(-2 * 24 * time.Hour), "2 days ago"},
		{"oneDay", time.Now().Add(-1 * 24 * time.Hour), "1 day ago"},
		{"weeks", time.Now().Add(-2 * 7 * 24 * time.Hour), "2 weeks ago"},
		{"oneWeek", time.Now().Add(-7 * 24 * time.Hour), "1 week ago"},
		{"months", time.Now().Add(-60 * 24 * time.Hour), "2 months ago"},
		{"years", time.Now().Add(-400 * 24 * time.Hour), "1 year ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{PushedAt: tt.pushedAt}
			if got := r.TimeSincePush(); got != tt.expected {
				t.Errorf("TimeSincePush() = %q, want %q", got, tt.expected)
			}
		})
	}
}
