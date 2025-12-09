package commitview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"gh-commit-analyzer/internal/models"
	"gh-commit-analyzer/internal/tui"
)

type Model struct {
	viewport     viewport.Model
	repoCommits  []models.RepoCommits
	expanded     map[int]bool
	cursor       int
	totalCommits int
	width        int
	height       int
	ready        bool
}

type RestartMsg struct{}

func New(repoCommits []models.RepoCommits, width, height int) Model {
	vp := viewport.New(width, height-4)
	vp.Style = tui.BoxStyle

	m := Model{
		viewport:    vp,
		repoCommits: repoCommits,
		expanded:    make(map[int]bool),
		width:       width,
		height:      height,
	}
	m.totalCommits = m.countCommits()
	m.updateContent()
	m.ready = true

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4
		m.updateContent()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tui.Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, tui.Keys.Restart):
			return m, func() tea.Msg { return RestartMsg{} }
		case key.Matches(msg, tui.Keys.Confirm):
			m.toggleExpanded()
			m.updateContent()
			return m, nil
		case key.Matches(msg, tui.Keys.Up):
			m.moveCursor(-1)
			m.updateContent()
			return m, nil
		case key.Matches(msg, tui.Keys.Down):
			m.moveCursor(1)
			m.updateContent()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	title := tui.TitleStyle.Render("Commits")
	help := tui.HelpStyle.Render("↑/↓: navigate • enter: expand • r: restart • q: quit")

	return fmt.Sprintf("%s\n%s\n%s", title, m.viewport.View(), help)
}

func (m *Model) updateContent() {
	var content strings.Builder
	commitIndex := 0

	for _, rc := range m.repoCommits {
		header := fmt.Sprintf("═══ %s (%s) - %d commits ═══",
			rc.Repository.NameWithOwner, rc.Branch, len(rc.Commits))
		content.WriteString(tui.RepoHeaderStyle.Render(header))
		content.WriteString("\n\n")

		for _, c := range rc.Commits {
			content.WriteString(m.renderCommit(c, commitIndex))
			content.WriteString("\n")
			commitIndex++
		}
		content.WriteString("\n")
	}

	m.viewport.SetContent(content.String())
}

func (m Model) renderCommit(c models.Commit, index int) string {
	cursor := "  "
	if index == m.cursor {
		cursor = "> "
	}

	sha := tui.CommitSHAStyle.Render(c.ShortSHA())
	date := tui.CommitDateStyle.Render(c.FormattedDate())
	author := tui.CommitAuthorStyle.Render(c.Author)

	header := fmt.Sprintf("%s%s │ %s │ %s", cursor, sha, date, author)

	if m.expanded[index] {
		return header + "\n" + m.renderExpandedMessage(c)
	}

	message := c.FirstLine()
	if c.HasMultipleLines() {
		message += tui.DimStyle.Render(fmt.Sprintf(" [+%d lines]", c.ExtraLineCount()))
	}

	return header + "\n" + cursor + "  └─ " + message
}

func (m Model) renderExpandedMessage(c models.Commit) string {
	var lines strings.Builder
	lines.WriteString("   ┌─────────────────────────────────────\n")
	lines.WriteString(fmt.Sprintf("   │ SHA:    %s\n", c.SHA))
	lines.WriteString(fmt.Sprintf("   │ Author: %s\n", c.AuthorWithEmail()))
	lines.WriteString(fmt.Sprintf("   │ Date:   %s\n", c.FormattedDate()))
	lines.WriteString("   │\n")

	for _, line := range strings.Split(c.Message, "\n") {
		lines.WriteString(fmt.Sprintf("   │ %s\n", line))
	}

	lines.WriteString("   └─────────────────────────────────────")
	return lines.String()
}

func (m *Model) toggleExpanded() {
	m.expanded[m.cursor] = !m.expanded[m.cursor]
}

func (m *Model) moveCursor(delta int) {
	m.cursor += delta
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= m.totalCommits {
		m.cursor = m.totalCommits - 1
	}
}

func (m Model) countCommits() int {
	count := 0
	for _, rc := range m.repoCommits {
		count += len(rc.Commits)
	}
	return count
}
