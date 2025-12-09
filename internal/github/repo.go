package github

import (
	"time"

	"gh-commit-analyzer/internal/models"
)

type repoResponse struct {
	Name             string `json:"name"`
	NameWithOwner    string `json:"nameWithOwner"`
	Description      string `json:"description"`
	URL              string `json:"url"`
	PushedAt         string `json:"pushedAt"`
	DefaultBranchRef struct {
		Name string `json:"name"`
	} `json:"defaultBranchRef"`
}

type branchResponse struct {
	Name string `json:"name"`
}

func ListRepositories() ([]models.Repository, error) {
	var response []repoResponse
	err := runGHWithJSON(&response,
		"repo", "list",
		"--json", "name,nameWithOwner,description,url,pushedAt,defaultBranchRef",
		"--limit", "1000",
	)
	if err != nil {
		return nil, err
	}
	return mapRepositories(response), nil
}

func ListBranches(owner, repo string) ([]string, error) {
	var response []branchResponse
	err := runGHWithJSON(&response, "api", "repos/"+owner+"/"+repo+"/branches")
	if err != nil {
		return nil, err
	}
	return extractBranchNames(response), nil
}

func mapRepositories(responses []repoResponse) []models.Repository {
	repos := make([]models.Repository, len(responses))
	for i, r := range responses {
		repos[i] = mapRepository(r)
	}
	return repos
}

func mapRepository(r repoResponse) models.Repository {
	pushedAt, _ := time.Parse(time.RFC3339, r.PushedAt)
	return models.Repository{
		Name:              r.Name,
		NameWithOwner:     r.NameWithOwner,
		Description:       r.Description,
		URL:               r.URL,
		PushedAt:          pushedAt,
		DefaultBranchName: r.DefaultBranchRef.Name,
	}
}

func extractBranchNames(responses []branchResponse) []string {
	names := make([]string, len(responses))
	for i, b := range responses {
		names[i] = b.Name
	}
	return names
}
