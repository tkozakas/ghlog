package reposelect

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gh-commit-analyzer/internal/models"
	"gh-commit-analyzer/internal/tui"
)

type item struct {
	repo     models.Repository
	selected bool
}

func (i item) Title() string {
	checkbox := "[ ]"
	if i.selected {
		checkbox = "[x]"
	}
	return fmt.Sprintf("%s %s", checkbox, i.repo.NameWithOwner)
}

func (i item) Description() string {
	desc := i.repo.TimeSincePush()
	if i.repo.Description != "" {
		desc = i.repo.Description + " • " + desc
	}
	return desc
}

func (i item) FilterValue() string {
	return i.repo.NameWithOwner
}

type Model struct {
	list     list.Model
	items    []item
	selected map[string]models.Repository
	width    int
	height   int
}

type DoneMsg struct {
	Selected []models.Repository
}

func New(repos []models.Repository, width, height int) Model {
	items := make([]item, len(repos))
	listItems := make([]list.Item, len(repos))
	for i, r := range repos {
		items[i] = item{repo: r}
		listItems[i] = items[i]
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(tui.ColorPrimary).
		BorderLeftForeground(tui.ColorPrimary)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(tui.ColorSecondary).
		BorderLeftForeground(tui.ColorPrimary)

	l := list.New(listItems, delegate, width, height-4)
	l.Title = "Select Repositories"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = tui.TitleStyle
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(tui.ColorPrimary)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(tui.ColorPrimary)
	l.SetStatusBarItemName("repo", "repos")
	l.AdditionalShortHelpKeys = shortHelpKeys
	l.AdditionalFullHelpKeys = fullHelpKeys

	return Model{
		list:     l,
		items:    items,
		selected: make(map[string]models.Repository),
		width:    width,
		height:   height,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, tui.Keys.Select):
			return m.toggleSelection(), nil
		case key.Matches(msg, tui.Keys.Confirm):
			if len(m.selected) > 0 {
				return m, m.confirmSelection
			}
		case key.Matches(msg, tui.Keys.Quit):
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	status := fmt.Sprintf("\n  %d selected", len(m.selected))
	if len(m.selected) > 0 {
		status = tui.SelectedStyle.Render(status)
	} else {
		status = tui.DimStyle.Render(status)
	}
	return m.list.View() + status
}

func (m Model) Selected() []models.Repository {
	repos := make([]models.Repository, 0, len(m.selected))
	for _, r := range m.selected {
		repos = append(repos, r)
	}
	return repos
}

func (m Model) toggleSelection() Model {
	if len(m.list.Items()) == 0 {
		return m
	}

	idx := m.list.Index()
	visibleItems := m.list.Items()
	if idx >= len(visibleItems) {
		return m
	}

	selectedItem := visibleItems[idx].(item)
	repoKey := selectedItem.repo.NameWithOwner

	if _, exists := m.selected[repoKey]; exists {
		delete(m.selected, repoKey)
	} else {
		m.selected[repoKey] = selectedItem.repo
	}

	m.updateListItems()
	return m
}

func (m *Model) updateListItems() {
	for i, it := range m.items {
		_, selected := m.selected[it.repo.NameWithOwner]
		m.items[i].selected = selected
	}

	listItems := make([]list.Item, len(m.items))
	for i, it := range m.items {
		listItems[i] = it
	}
	m.list.SetItems(listItems)
}

func (m Model) confirmSelection() tea.Msg {
	return DoneMsg{Selected: m.Selected()}
}

func shortHelpKeys() []key.Binding {
	return []key.Binding{tui.Keys.Select, tui.Keys.Confirm}
}

func fullHelpKeys() []key.Binding {
	return []key.Binding{tui.Keys.Select, tui.Keys.Confirm, tui.Keys.Quit}
}
