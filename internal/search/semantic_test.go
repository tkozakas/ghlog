package search

import (
	"testing"

	"github.com/tkozakas/gh-log/internal/models"
)

func TestExtractSHAFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"fullPath", "/tmp/gh-commit-semantic-123/abc123def456.txt", "abc123def456"},
		{"shortPath", "/tmp/abc.txt", "abc"},
		{"fileOnly", "simple.txt", "simple"},
		{"nestedPath", "/path/to/sha123.txt", "sha123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractSHAFromPath(tt.path); got != tt.expected {
				t.Errorf("extractSHAFromPath() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFilterCommitsSemantically_EmptyQuery(t *testing.T) {
	commits := []models.Commit{
		{SHA: "abc123", Message: "Fix bug"},
		{SHA: "def456", Message: "Add feature"},
	}

	result, err := FilterCommitsSemantically(commits, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(commits) {
		t.Errorf("len(result) = %d, want %d", len(result), len(commits))
	}
}

func TestFilterCommitsSemantically_EmptyCommits(t *testing.T) {
	result, err := FilterCommitsSemantically([]models.Commit{}, "bug fix")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("len(result) = %d, want 0", len(result))
	}
}

func TestFilterAndSortByScore(t *testing.T) {
	commits := []models.Commit{
		{SHA: "aaa", Message: "First"},
		{SHA: "bbb", Message: "Second"},
		{SHA: "ccc", Message: "Third"},
	}
	scores := map[string]float64{
		"aaa": 0.5,
		"ccc": 0.9,
	}

	result := filterAndSortByScore(commits, scores)

	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}
	if result[0].SHA != "ccc" {
		t.Errorf("result[0].SHA = %q, want %q", result[0].SHA, "ccc")
	}
	if result[1].SHA != "aaa" {
		t.Errorf("result[1].SHA = %q, want %q", result[1].SHA, "aaa")
	}
}

func TestFilterAndSortByScore_EmptyScores(t *testing.T) {
	commits := []models.Commit{{SHA: "aaa", Message: "First"}}

	result := filterAndSortByScore(commits, map[string]float64{})

	if len(result) != 0 {
		t.Errorf("len(result) = %d, want 0", len(result))
	}
}

func TestParseCkOutput(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{
			"validLines",
			`{"file":"/tmp/abc123.txt","score":0.85}
{"file":"/tmp/def456.txt","score":0.72}
`,
			2,
		},
		{
			"emptyLines",
			`{"file":"/tmp/abc123.txt","score":0.85}

{"file":"/tmp/def456.txt","score":0.72}
`,
			2,
		},
		{
			"invalidJSON",
			`{"file":"/tmp/abc123.txt","score":0.85}
invalid json line
{"file":"/tmp/def456.txt","score":0.72}
`,
			2,
		},
		{
			"emptyInput",
			"",
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseCkOutput([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != tt.expectedCount {
				t.Errorf("len(result) = %d, want %d", len(result), tt.expectedCount)
			}
		})
	}
}

func TestParseCkOutput_ScoreValues(t *testing.T) {
	input := `{"file":"/tmp/abc123.txt","score":0.85}
{"file":"/tmp/def456.txt","score":0.72}
`
	result, err := parseCkOutput([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["abc123"] != 0.85 {
		t.Errorf("score[abc123] = %f, want 0.85", result["abc123"])
	}
	if result["def456"] != 0.72 {
		t.Errorf("score[def456] = %f, want 0.72", result["def456"])
	}
}
