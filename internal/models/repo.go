package models

import (
	"fmt"
	"strings"
	"time"
)

type Repository struct {
	Name              string    `json:"name"`
	NameWithOwner     string    `json:"nameWithOwner"`
	Description       string    `json:"description"`
	URL               string    `json:"url"`
	PushedAt          time.Time `json:"pushedAt"`
	DefaultBranchName string    `json:"defaultBranchName"`
}

func (r Repository) Owner() string {
	return splitNameWithOwner(r.NameWithOwner, 0)
}

func (r Repository) RepoName() string {
	if part := splitNameWithOwner(r.NameWithOwner, 1); part != "" {
		return part
	}
	return r.Name
}

func (r Repository) TimeSincePush() string {
	if r.PushedAt.IsZero() {
		return "unknown"
	}
	return formatDuration(time.Since(r.PushedAt))
}

func splitNameWithOwner(nameWithOwner string, index int) string {
	parts := strings.SplitN(nameWithOwner, "/", 2)
	if index < len(parts) {
		return parts[index]
	}
	return ""
}

func formatDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return pluralize(int(d.Minutes()), "minute")
	case d < 24*time.Hour:
		return pluralize(int(d.Hours()), "hour")
	case d < 7*24*time.Hour:
		return pluralize(int(d.Hours()/24), "day")
	case d < 30*24*time.Hour:
		return pluralize(int(d.Hours()/24/7), "week")
	case d < 365*24*time.Hour:
		return pluralize(int(d.Hours()/24/30), "month")
	default:
		return pluralize(int(d.Hours()/24/365), "year")
	}
}

func pluralize(count int, unit string) string {
	if count == 1 {
		return fmt.Sprintf("1 %s ago", unit)
	}
	return fmt.Sprintf("%d %ss ago", count, unit)
}
