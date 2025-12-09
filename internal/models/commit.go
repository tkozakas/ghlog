package models

import (
	"fmt"
	"strings"
	"time"
)

type Commit struct {
	SHA     string    `json:"sha"`
	Message string    `json:"message"`
	Author  string    `json:"author"`
	Email   string    `json:"email"`
	Date    time.Time `json:"date"`
	URL     string    `json:"url"`
}

type RepoCommits struct {
	Repository Repository
	Branch     string
	Commits    []Commit
	HasMore    bool
	Page       int
	TotalCount int
}

func (c Commit) ShortSHA() string {
	if len(c.SHA) >= 7 {
		return c.SHA[:7]
	}
	return c.SHA
}

func (c Commit) FirstLine() string {
	return strings.TrimSpace(strings.SplitN(c.Message, "\n", 2)[0])
}

func (c Commit) HasMultipleLines() bool {
	return strings.Contains(strings.TrimSpace(c.Message), "\n")
}

func (c Commit) ExtraLineCount() int {
	lines := strings.Split(strings.TrimSpace(c.Message), "\n")
	if len(lines) > 1 {
		return len(lines) - 1
	}
	return 0
}

func (c Commit) FormattedDate() string {
	if c.Date.IsZero() {
		return "unknown"
	}
	return c.Date.Format("2006-01-02 15:04")
}

func (c Commit) AuthorWithEmail() string {
	if c.Email != "" {
		return fmt.Sprintf("%s <%s>", c.Author, c.Email)
	}
	return c.Author
}
