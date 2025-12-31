package branchselect

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"ghlog/internal/models"
	"ghlog/internal/tui"
)

type item struct {
	name      string
	isDefault bool
}

func (i item) Title() string {
	if i.isDefault {
		return i.name + " (default)"
	}
	return i.name
}

func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.name }

type Model struct {
	list           list.Model
	repo           models.Repository
	selectedBranch string
}

type DoneMsg struct {
	Repo   models.Repository
	Branch string
}

type UseDefaultMsg struct{}

func New(repo models.Repository, branches []string, width, height int) Model {
	defaultBranch := repo.DefaultBranchName
	items := make([]list.Item, len(branches))

	for i, b := range branches {
		items[i] = item{name: b, isDefault: b == defaultBranch}
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(tui.ColorPrimary).
		BorderLeftForeground(tui.ColorPrimary)

	l := list.New(items, delegate, width, height-4)
	l.Title = fmt.Sprintf("Select branch for %s", repo.NameWithOwner)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = tui.TitleStyle
	l.AdditionalShortHelpKeys = shortHelpKeys

	return Model{
		list: l,
		repo: repo,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, tui.Keys.Confirm):
			if len(m.list.Items()) > 0 {
				selected := m.list.SelectedItem().(item)
				m.selectedBranch = selected.name
				return m, m.submit
			}
		case key.Matches(msg, tui.Keys.Default):
			return m, m.submitDefault
		case key.Matches(msg, tui.Keys.Back):
			m.selectedBranch = m.repo.DefaultBranchName
			return m, m.submit
		case key.Matches(msg, tui.Keys.Quit):
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	help := tui.HelpStyle.Render("  d: use default for all remaining")
	return m.list.View() + "\n" + help
}

func (m Model) submit() tea.Msg {
	branch := m.selectedBranch
	if branch == "" {
		branch = m.repo.DefaultBranchName
	}
	return DoneMsg{Repo: m.repo, Branch: branch}
}

func (m Model) submitDefault() tea.Msg {
	return UseDefaultMsg{}
}

func shortHelpKeys() []key.Binding {
	return []key.Binding{tui.Keys.Default}
}
