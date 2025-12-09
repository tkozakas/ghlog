package filterform

import (
	"fmt"
	"strconv"

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
	fieldCount
)

type Model struct {
	inputs  []textinput.Model
	focused int
}

type DoneMsg struct {
	Filters models.FilterOptions
}

func New() Model {
	inputs := make([]textinput.Model, fieldCount)

	inputs[fieldDateFrom] = newInput("2024-01-01", 10)
	inputs[fieldDateTo] = newInput("2024-12-31", 10)
	inputs[fieldAuthor] = newInput("username", 30)
	inputs[fieldPerPage] = newInput("50", 3)
	inputs[fieldPerPage].SetValue("50")

	inputs[fieldDateFrom].Focus()

	return Model{inputs: inputs}
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
		case key.Matches(msg, tui.Keys.Back):
			return m, m.submit
		case key.Matches(msg, tui.Keys.Confirm):
			return m, m.submit
		case key.Matches(msg, tui.Keys.Tab):
			return m.nextField(), nil
		case key.Matches(msg, tui.Keys.ShiftTab):
			return m.prevField(), nil
		}
	}

	return m.updateInputs(msg)
}

func (m Model) View() string {
	title := tui.TitleStyle.Render("Configure Filters")
	help := tui.HelpStyle.Render("tab: next • shift+tab: prev • enter: confirm • esc: skip")

	form := fmt.Sprintf(
		"%s\n\n  %s  %s\n\n  %s\n\n  %s\n\n%s",
		title,
		m.renderField(fieldDateFrom, "From:    "),
		m.renderField(fieldDateTo, "To: "),
		m.renderField(fieldAuthor, "Author:  "),
		m.renderField(fieldPerPage, "Per page:"),
		help,
	)

	return form
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

func (m Model) nextField() Model {
	m.inputs[m.focused].Blur()
	m.focused = (m.focused + 1) % fieldCount
	m.inputs[m.focused].Focus()
	return m
}

func (m Model) prevField() Model {
	m.inputs[m.focused].Blur()
	m.focused--
	if m.focused < 0 {
		m.focused = fieldCount - 1
	}
	m.inputs[m.focused].Focus()
	return m
}

func (m Model) updateInputs(msg tea.Msg) (Model, tea.Cmd) {
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

func (m Model) submit() tea.Msg {
	return DoneMsg{Filters: m.Filters()}
}

func newInput(placeholder string, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Width = width
	ti.CharLimit = width + 10
	return ti
}
