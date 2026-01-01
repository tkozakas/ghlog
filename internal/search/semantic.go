package search

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tkozakas/gh-log/internal/models"
)

type ckResult struct {
	File  string  `json:"file"`
	Score float64 `json:"score"`
}

type scoredCommit struct {
	commit models.Commit
	score  float64
}

func IsAvailable() bool {
	_, err := exec.LookPath("ck")
	return err == nil
}

func FilterCommitsSemantically(commits []models.Commit, query string) ([]models.Commit, error) {
	if query == "" || len(commits) == 0 {
		return commits, nil
	}

	tmpDir, err := createTempCommitFiles(commits)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp files: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	matchedSHAs, err := runCkSearch(tmpDir, query)
	if err != nil {
		return nil, err
	}

	return filterAndSortByScore(commits, matchedSHAs), nil
}

func createTempCommitFiles(commits []models.Commit) (string, error) {
	tmpDir, err := os.MkdirTemp("", "gh-commit-semantic-*")
	if err != nil {
		return "", err
	}

	for _, c := range commits {
		filePath := filepath.Join(tmpDir, c.SHA+".txt")
		if err := os.WriteFile(filePath, []byte(c.Message), 0644); err != nil {
			os.RemoveAll(tmpDir)
			return "", err
		}
	}

	return tmpDir, nil
}

func runCkSearch(dir, query string) (map[string]float64, error) {
	if err := indexDirectory(dir); err != nil {
		return nil, err
	}

	cmd := exec.Command("ck", "--jsonl", "--hybrid", "--scores", query, dir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if isNoMatchesError(err, stderr.Bytes()) {
			return make(map[string]float64), nil
		}
		return nil, fmt.Errorf("ck search failed: %w: %s", err, stderr.String())
	}

	return parseCkOutput(stdout.Bytes())
}

func indexDirectory(dir string) error {
	cmd := exec.Command("ck", "--index", dir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ck index failed: %w: %s", err, stderr.String())
	}
	return nil
}

func isNoMatchesError(err error, stderr []byte) bool {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode() == 1 && len(stderr) == 0
	}
	return false
}

func parseCkOutput(output []byte) (map[string]float64, error) {
	result := make(map[string]float64)
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var ckRes ckResult
		if err := json.Unmarshal([]byte(line), &ckRes); err != nil {
			continue
		}

		sha := extractSHAFromPath(ckRes.File)
		if sha != "" {
			result[sha] = ckRes.Score
		}
	}

	return result, scanner.Err()
}

func extractSHAFromPath(filePath string) string {
	base := filepath.Base(filePath)
	return strings.TrimSuffix(base, ".txt")
}

func filterAndSortByScore(commits []models.Commit, scores map[string]float64) []models.Commit {
	if len(scores) == 0 {
		return []models.Commit{}
	}

	var scored []scoredCommit
	for _, c := range commits {
		if score, ok := scores[c.SHA]; ok {
			scored = append(scored, scoredCommit{commit: c, score: score})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	result := make([]models.Commit, len(scored))
	for i, sc := range scored {
		result[i] = sc.commit
	}

	return result
}
