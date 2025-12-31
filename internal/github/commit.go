package github

import (
	"fmt"
	"time"

	"ghlog/internal/models"
)

type commitResponse struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Date  string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
	HTMLURL string `json:"html_url"`
}

func GetCommits(owner, repo, branch string, filters models.FilterOptions, page int) ([]models.Commit, bool, error) {
	endpoint := buildCommitsEndpoint(owner, repo, branch, filters, page)

	var response []commitResponse
	if err := runGHWithJSON(&response, "api", endpoint); err != nil {
		return nil, false, err
	}

	commits := mapCommits(response)
	hasMore := len(commits) == filters.PerPage
	return commits, hasMore, nil
}

func buildCommitsEndpoint(owner, repo, branch string, filters models.FilterOptions, page int) string {
	endpoint := fmt.Sprintf("repos/%s/%s/commits?per_page=%d&page=%d",
		owner, repo, filters.PerPage, page)

	if branch != "" {
		endpoint += "&sha=" + branch
	}
	if filters.DateFrom != "" {
		endpoint += "&since=" + filters.DateFrom + "T00:00:00Z"
	}
	if filters.DateTo != "" {
		endpoint += "&until=" + filters.DateTo + "T23:59:59Z"
	}
	if filters.Author != "" {
		endpoint += "&author=" + filters.Author
	}
	return endpoint
}

func mapCommits(responses []commitResponse) []models.Commit {
	commits := make([]models.Commit, len(responses))
	for i, r := range responses {
		commits[i] = mapCommit(r)
	}
	return commits
}

func mapCommit(r commitResponse) models.Commit {
	date, _ := time.Parse(time.RFC3339, r.Commit.Author.Date)
	return models.Commit{
		SHA:     r.SHA,
		Message: r.Commit.Message,
		Author:  r.Commit.Author.Name,
		Email:   r.Commit.Author.Email,
		Date:    date,
		URL:     r.HTMLURL,
	}
}
