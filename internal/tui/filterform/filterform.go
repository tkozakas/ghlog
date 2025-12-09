package filterform

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"gh-commit-analyzer/internal/models"
	"gh-commit-analyzer/internal/tui"
)

const (
	fieldDateFrom = iota
	fieldDateTo
	fieldAuthor
	fieldPerPage
	fieldCountBase
)

type RepoBranches struct {
	Repo     models.Repository
	Branches []string
}

type Model struct {
	inputs       []textinput.Model
	repoBranches []RepoBranches
	branchIdx    []int
	focused      int
	fieldCount   int
}

type DoneMsg struct {
	Filters  models.FilterOptions
	Branches map[string]string
}

func New(repoBranches []RepoBranches) Model {
	fieldCount := fieldCountBase + len(repoBranches)
	inputs := make([]textinput.Model, fieldCountBase)

	inputs[fieldDateFrom] = newInput("2024-01-01", 10)
	inputs[fieldDateTo] = newInput("2024-12-31", 10)
	inputs[fieldAuthor] = newInput("username", 30)
	inputs[fieldPerPage] = newInput("50", 3)
	inputs[fieldPerPage].SetValue("50")

	inputs[fieldDateFrom].Focus()

	branchIdx := make([]int, len(repoBranches))
	for i, rb := range repoBranches {
		branchIdx[i] = findDefaultBranchIndex(rb)
	}

	return Model{
		inputs:       inputs,
		repoBranches: repoBranches,
		branchIdx:    branchIdx,
		fieldCount:   fieldCount,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tui.Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, tui.Keys.Back), key.Matches(msg, tui.Keys.Confirm):
			return m, m.submit
		case key.Matches(msg, tui.Keys.Tab):
			return m.nextField(), nil
		case key.Matches(msg, tui.Keys.ShiftTab):
			return m.prevField(), nil
		case key.Matches(msg, tui.Keys.Up), key.Matches(msg, tui.Keys.Down):
			if m.isBranchField() {
				return m.cycleBranch(msg), nil
			}
		}
	}

	return m.updateInputs(msg)
}

func (m Model) View() string {
	title := tui.TitleStyle.Render("Configure Filters")
	help := tui.HelpStyle.Render("tab: next • shift+tab: prev • ↑/↓: cycle branch • enter: confirm")

	var b strings.Builder
	b.WriteString(title)
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("  %s  %s\n\n",
		m.renderField(fieldDateFrom, "From:    "),
		m.renderField(fieldDateTo, "To: ")))
	b.WriteString(fmt.Sprintf("  %s\n\n", m.renderField(fieldAuthor, "Author:  ")))
	b.WriteString(fmt.Sprintf("  %s\n\n", m.renderField(fieldPerPage, "Per page:")))

	if len(m.repoBranches) > 0 {
		b.WriteString("  " + tui.DimStyle.Render("─── Branches ───") + "\n\n")
		for i, rb := range m.repoBranches {
			b.WriteString(fmt.Sprintf("  %s\n", m.renderBranchField(i, rb)))
		}
	}

	b.WriteString("\n")
	b.WriteString(help)

	return b.String()
}

func (m Model) Filters() models.FilterOptions {
	perPage, _ := strconv.Atoi(m.inputs[fieldPerPage].Value())
	filters := models.FilterOptions{
		DateFrom: m.inputs[fieldDateFrom].Value(),
		DateTo:   m.inputs[fieldDateTo].Value(),
		Author:   m.inputs[fieldAuthor].Value(),
		PerPage:  perPage,
	}
	filters.Validate()
	return filters
}

func (m Model) Branches() map[string]string {
	result := make(map[string]string)
	for i, rb := range m.repoBranches {
		idx := m.branchIdx[i]
		if idx >= 0 && idx < len(rb.Branches) {
			result[rb.Repo.NameWithOwner] = rb.Branches[idx]
		} else {
			result[rb.Repo.NameWithOwner] = rb.Repo.DefaultBranchName
		}
	}
	return result
}

func (m Model) isBranchField() bool {
	return m.focused >= fieldCountBase
}

func (m Model) cycleBranch(msg tea.KeyMsg) Model {
	repoIdx := m.focused - fieldCountBase
	if repoIdx < 0 || repoIdx >= len(m.repoBranches) {
		return m
	}

	branchCount := len(m.repoBranches[repoIdx].Branches)
	if branchCount == 0 {
		return m
	}

	if key.Matches(msg, tui.Keys.Down) {
		m.branchIdx[repoIdx] = (m.branchIdx[repoIdx] + 1) % branchCount
	} else {
		m.branchIdx[repoIdx]--
		if m.branchIdx[repoIdx] < 0 {
			m.branchIdx[repoIdx] = branchCount - 1
		}
	}
	return m
}

func (m Model) nextField() Model {
	m.blurCurrent()
	m.focused = (m.focused + 1) % m.fieldCount
	m.focusCurrent()
	return m
}

func (m Model) prevField() Model {
	m.blurCurrent()
	m.focused--
	if m.focused < 0 {
		m.focused = m.fieldCount - 1
	}
	m.focusCurrent()
	return m
}

func (m *Model) blurCurrent() {
	if m.focused < fieldCountBase {
		m.inputs[m.focused].Blur()
	}
}

func (m *Model) focusCurrent() {
	if m.focused < fieldCountBase {
		m.inputs[m.focused].Focus()
	}
}

func (m Model) updateInputs(msg tea.Msg) (Model, tea.Cmd) {
	if m.focused >= fieldCountBase {
		return m, nil
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) renderField(field int, label string) string {
	style := tui.DimStyle
	if m.focused == field {
		style = tui.SelectedStyle
	}
	return style.Render(label) + m.inputs[field].View()
}

func (m Model) renderBranchField(repoIdx int, rb RepoBranches) string {
	fieldIdx := fieldCountBase + repoIdx
	isFocused := m.focused == fieldIdx

	labelStyle := tui.DimStyle
	if isFocused {
		labelStyle = tui.SelectedStyle
	}

	branchName := rb.Repo.DefaultBranchName
	if m.branchIdx[repoIdx] >= 0 && m.branchIdx[repoIdx] < len(rb.Branches) {
		branchName = rb.Branches[m.branchIdx[repoIdx]]
	}

	isDefault := branchName == rb.Repo.DefaultBranchName
	branchDisplay := branchName
	if isDefault {
		branchDisplay += " (default)"
	}

	cursor := "  "
	if isFocused {
		cursor = "> "
	}

	return fmt.Sprintf("%s%s %s",
		cursor,
		labelStyle.Render(rb.Repo.NameWithOwner+":"),
		tui.CommitSHAStyle.Render(branchDisplay))
}

func (m Model) submit() tea.Msg {
	return DoneMsg{
		Filters:  m.Filters(),
		Branches: m.Branches(),
	}
}

func findDefaultBranchIndex(rb RepoBranches) int {
	for i, b := range rb.Branches {
		if b == rb.Repo.DefaultBranchName {
			return i
		}
	}
	return 0
}

func newInput(placeholder string, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Width = width
	ti.CharLimit = width + 10
	return ti
}
