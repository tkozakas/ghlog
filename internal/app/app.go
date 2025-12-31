package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"ghlog/internal/github"
	"ghlog/internal/models"
	"ghlog/internal/search"
	"ghlog/internal/tui"
	"ghlog/internal/tui/commitview"
	"ghlog/internal/tui/filterform"
	"ghlog/internal/tui/reposelect"
)

type state int

const (
	stateLoading state = iota
	stateRepoSelect
	stateLoadingBranches
	stateFilterForm
	stateLoadingCommits
	stateCommitView
	stateError
)

type Model struct {
	state         state
	width         int
	height        int
	spinner       spinner.Model
	err           error
	repos         []models.Repository
	selectedRepos []models.Repository
	repoBranches  []filterform.RepoBranches
	filters       models.FilterOptions
	branches      map[string]string
	repoCommits   []models.RepoCommits
	repoSelect    reposelect.Model
	filterForm    filterform.Model
	commitView    commitview.Model
}

type reposLoadedMsg struct{ repos []models.Repository }
type allBranchesLoadedMsg struct{ repoBranches []filterform.RepoBranches }
type commitsLoadedMsg struct{ repoCommits []models.RepoCommits }
type moreCommitsLoadedMsg struct {
	repoName string
	commits  []models.Commit
	page     int
	hasMore  bool
}
type errMsg struct{ err error }

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = tui.SelectedStyle

	return Model{
		state:    stateLoading,
		spinner:  s,
		branches: make(map[string]string),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, loadRepos)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m.propagateSize(), nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case reposLoadedMsg:
		m.repos = msg.repos
		m.repoSelect = reposelect.New(m.repos, m.width, m.height)
		m.state = stateRepoSelect
		return m, nil

	case reposelect.DoneMsg:
		m.selectedRepos = msg.Selected
		m.state = stateLoadingBranches
		return m, m.loadAllBranches()

	case allBranchesLoadedMsg:
		m.repoBranches = msg.repoBranches
		m.filterForm = filterform.New(m.repoBranches)
		m.state = stateFilterForm
		return m, nil

	case filterform.DoneMsg:
		m.filters = msg.Filters
		m.branches = msg.Branches
		return m.loadAllCommits()

	case commitsLoadedMsg:
		m.repoCommits = msg.repoCommits
		m.commitView = commitview.New(m.repoCommits, m.width, m.height)
		m.state = stateCommitView
		return m, nil

	case moreCommitsLoadedMsg:
		m.repoCommits = updateRepoCommits(m.repoCommits, msg)
		m.commitView.UpdateCommits(m.repoCommits)
		return m, nil

	case commitview.LoadMoreMsg:
		return m, m.loadMoreCommits(msg.RepoName, msg.NextPage)

	case commitview.RestartMsg:
		return m.restart()

	case errMsg:
		m.err = msg.err
		m.state = stateError
		return m, nil
	}

	return m.updateChild(msg)
}

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return m.viewLoading("Loading repositories...")
	case stateRepoSelect:
		return m.repoSelect.View()
	case stateLoadingBranches:
		return m.viewLoading("Loading branches...")
	case stateFilterForm:
		return m.filterForm.View()
	case stateLoadingCommits:
		return m.viewLoading("Loading commits...")
	case stateCommitView:
		return m.commitView.View()
	case stateError:
		return m.viewError()
	default:
		return ""
	}
}

func (m Model) updateChild(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case stateRepoSelect:
		m.repoSelect, cmd = m.repoSelect.Update(msg)
	case stateFilterForm:
		m.filterForm, cmd = m.filterForm.Update(msg)
	case stateCommitView:
		m.commitView, cmd = m.commitView.Update(msg)
	}

	return m, cmd
}

func (m Model) propagateSize() Model {
	switch m.state {
	case stateRepoSelect:
		m.repoSelect, _ = m.repoSelect.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	case stateCommitView:
		m.commitView, _ = m.commitView.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	}
	return m
}

func (m Model) loadAllBranches() tea.Cmd {
	return func() tea.Msg {
		var repoBranches []filterform.RepoBranches

		for _, repo := range m.selectedRepos {
			branches, err := github.ListBranches(repo.Owner(), repo.RepoName())
			if err != nil {
				return errMsg{err: err}
			}
			repoBranches = append(repoBranches, filterform.RepoBranches{
				Repo:     repo,
				Branches: branches,
			})
		}

		return allBranchesLoadedMsg{repoBranches: repoBranches}
	}
}

func (m Model) loadAllCommits() (Model, tea.Cmd) {
	m.state = stateLoadingCommits
	return m, m.loadCommitsCmd()
}

func (m Model) loadCommitsCmd() tea.Cmd {
	return func() tea.Msg {
		var repoCommits []models.RepoCommits
		for _, repo := range m.selectedRepos {
			branch := m.branches[repo.NameWithOwner]
			if branch == "" {
				branch = repo.DefaultBranchName
			}

			commits, hasMore, err := github.GetCommits(repo.Owner(), repo.RepoName(), branch, m.filters, 1)
			if err != nil {
				return errMsg{err: err}
			}

			commits, err = applySemanticFilter(commits, m.filters.SemanticQuery)
			if err != nil {
				return errMsg{err: err}
			}

			repoCommits = append(repoCommits, models.RepoCommits{
				Repository: repo,
				Branch:     branch,
				Commits:    commits,
				Page:       1,
				HasMore:    hasMore,
			})
		}

		return commitsLoadedMsg{repoCommits: repoCommits}
	}
}

func (m Model) loadMoreCommits(repoName string, page int) tea.Cmd {
	return func() tea.Msg {
		for _, repo := range m.selectedRepos {
			if repo.NameWithOwner != repoName {
				continue
			}

			branch := m.branches[repo.NameWithOwner]
			if branch == "" {
				branch = repo.DefaultBranchName
			}

			commits, hasMore, err := github.GetCommits(repo.Owner(), repo.RepoName(), branch, m.filters, page)
			if err != nil {
				return errMsg{err: err}
			}

			return moreCommitsLoadedMsg{
				repoName: repoName,
				commits:  commits,
				page:     page,
				hasMore:  hasMore,
			}
		}
		return nil
	}
}

func (m Model) restart() (Model, tea.Cmd) {
	m.selectedRepos = nil
	m.branches = make(map[string]string)
	m.repoCommits = nil
	m.repoSelect = reposelect.New(m.repos, m.width, m.height)
	m.state = stateRepoSelect
	return m, nil
}

func (m Model) viewLoading(msg string) string {
	return fmt.Sprintf("\n  %s %s\n", m.spinner.View(), msg)
}

func (m Model) viewError() string {
	return tui.ErrorStyle.Render(fmt.Sprintf("\n  Error: %v\n\n  Press q to quit.\n", m.err))
}

func updateRepoCommits(repoCommits []models.RepoCommits, msg moreCommitsLoadedMsg) []models.RepoCommits {
	for i, rc := range repoCommits {
		if rc.Repository.NameWithOwner == msg.repoName {
			repoCommits[i].Commits = append(rc.Commits, msg.commits...)
			repoCommits[i].Page = msg.page
			repoCommits[i].HasMore = msg.hasMore
			break
		}
	}
	return repoCommits
}

func loadRepos() tea.Msg {
	repos, err := github.ListRepositories()
	if err != nil {
		return errMsg{err: err}
	}
	return reposLoadedMsg{repos: repos}
}

func applySemanticFilter(commits []models.Commit, query string) ([]models.Commit, error) {
	if query == "" {
		return commits, nil
	}
	if !search.IsAvailable() {
		return commits, nil
	}
	return search.FilterCommitsSemantically(commits, query)
}
