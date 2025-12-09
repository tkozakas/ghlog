package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"gh-commit-analyzer/internal/github"
	"gh-commit-analyzer/internal/models"
	"gh-commit-analyzer/internal/tui"
	"gh-commit-analyzer/internal/tui/branchselect"
	"gh-commit-analyzer/internal/tui/commitview"
	"gh-commit-analyzer/internal/tui/filterform"
	"gh-commit-analyzer/internal/tui/reposelect"
)

type state int

const (
	stateLoading state = iota
	stateRepoSelect
	stateFilterForm
	stateBranchSelect
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
	filters       models.FilterOptions
	branches      map[string]string
	branchIndex   int
	repoCommits   []models.RepoCommits
	repoSelect    reposelect.Model
	filterForm    filterform.Model
	branchSelect  branchselect.Model
	commitView    commitview.Model
}

type reposLoadedMsg struct{ repos []models.Repository }
type branchesLoadedMsg struct{ branches []string }
type commitsLoadedMsg struct{ repoCommits []models.RepoCommits }
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
		m.filterForm = filterform.New()
		m.state = stateFilterForm
		return m, nil

	case filterform.DoneMsg:
		m.filters = msg.Filters
		return m.startBranchSelection()

	case branchesLoadedMsg:
		repo := m.selectedRepos[m.branchIndex]
		m.branchSelect = branchselect.New(repo, msg.branches, m.width, m.height)
		m.state = stateBranchSelect
		return m, nil

	case branchselect.DoneMsg:
		m.branches[msg.Repo.NameWithOwner] = msg.Branch
		return m.nextBranchOrLoadCommits()

	case branchselect.UseDefaultMsg:
		return m.useDefaultBranches()

	case commitsLoadedMsg:
		m.repoCommits = msg.repoCommits
		m.commitView = commitview.New(m.repoCommits, m.width, m.height)
		m.state = stateCommitView
		return m, nil

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
	case stateFilterForm:
		return m.filterForm.View()
	case stateBranchSelect:
		return m.branchSelect.View()
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
	case stateBranchSelect:
		m.branchSelect, cmd = m.branchSelect.Update(msg)
	case stateCommitView:
		m.commitView, cmd = m.commitView.Update(msg)
	}

	return m, cmd
}

func (m Model) propagateSize() Model {
	switch m.state {
	case stateRepoSelect:
		m.repoSelect, _ = m.repoSelect.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	case stateBranchSelect:
		m.branchSelect, _ = m.branchSelect.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	case stateCommitView:
		m.commitView, _ = m.commitView.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	}
	return m
}

func (m Model) startBranchSelection() (Model, tea.Cmd) {
	m.branchIndex = 0
	return m.loadBranchesForCurrent()
}

func (m Model) loadBranchesForCurrent() (Model, tea.Cmd) {
	if m.branchIndex >= len(m.selectedRepos) {
		return m.loadAllCommits()
	}

	m.state = stateLoading
	repo := m.selectedRepos[m.branchIndex]
	return m, loadBranches(repo)
}

func (m Model) nextBranchOrLoadCommits() (Model, tea.Cmd) {
	m.branchIndex++
	return m.loadBranchesForCurrent()
}

func (m Model) useDefaultBranches() (Model, tea.Cmd) {
	for i := m.branchIndex; i < len(m.selectedRepos); i++ {
		repo := m.selectedRepos[i]
		m.branches[repo.NameWithOwner] = repo.DefaultBranchName
	}
	return m.loadAllCommits()
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

			commits, _, err := github.GetCommits(repo.Owner(), repo.RepoName(), branch, m.filters, 1)
			if err != nil {
				return errMsg{err: err}
			}

			repoCommits = append(repoCommits, models.RepoCommits{
				Repository: repo,
				Branch:     branch,
				Commits:    commits,
				Page:       1,
			})
		}

		return commitsLoadedMsg{repoCommits: repoCommits}
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

func loadRepos() tea.Msg {
	repos, err := github.ListRepositories()
	if err != nil {
		return errMsg{err: err}
	}
	return reposLoadedMsg{repos: repos}
}

func loadBranches(repo models.Repository) tea.Cmd {
	return func() tea.Msg {
		branches, err := github.ListBranches(repo.Owner(), repo.RepoName())
		if err != nil {
			return errMsg{err: err}
		}
		return branchesLoadedMsg{branches: branches}
	}
}
