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
	loading      bool
}

type RestartMsg struct{}

type LoadMoreMsg struct {
	RepoName string
	NextPage int
}

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

func (m *Model) UpdateCommits(repoCommits []models.RepoCommits) {
	m.repoCommits = repoCommits
	m.totalCommits = m.countCommits()
	m.loading = false
	m.updateContent()
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
			cmd := m.checkLoadMore()
			m.updateContent()
			return m, cmd
		case key.Matches(msg, tui.Keys.NextPage):
			cmd := m.loadMoreForCurrentRepo()
			m.updateContent()
			return m, cmd
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
	help := tui.HelpStyle.Render("↑/↓: navigate • enter: expand • n: load more • r: restart • q: quit")

	return fmt.Sprintf("%s\n%s\n%s", title, m.viewport.View(), help)
}

func (m *Model) updateContent() {
	var content strings.Builder
	commitIndex := 0
	cursorLine := 0
	lineCount := 0

	for _, rc := range m.repoCommits {
		header := fmt.Sprintf("═══ %s (%s) - %d commits ═══",
			rc.Repository.NameWithOwner, rc.Branch, len(rc.Commits))
		content.WriteString(tui.RepoHeaderStyle.Render(header))
		content.WriteString("\n\n")
		lineCount += 2

		for _, c := range rc.Commits {
			if commitIndex == m.cursor {
				cursorLine = lineCount
			}
			commitStr := m.renderCommit(c, commitIndex)
			content.WriteString(commitStr)
			content.WriteString("\n")
			lineCount += strings.Count(commitStr, "\n") + 1
			commitIndex++
		}

		if rc.HasMore {
			content.WriteString(tui.DimStyle.Render("    ↓ press 'n' to load more..."))
			content.WriteString("\n")
			lineCount += 1
		}
		content.WriteString("\n")
		lineCount += 1
	}

	if m.loading {
		content.WriteString(tui.SelectedStyle.Render("  Loading more commits..."))
		content.WriteString("\n")
	}

	m.viewport.SetContent(content.String())
	m.ensureCursorVisible(cursorLine)
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

	return header + "\n     └─ " + message
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
	maxCursor := m.totalCommits - 1
	if maxCursor < 0 {
		maxCursor = 0
	}
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}
}

func (m *Model) ensureCursorVisible(cursorLine int) {
	viewTop := m.viewport.YOffset
	viewBottom := viewTop + m.viewport.Height

	if cursorLine < viewTop {
		m.viewport.SetYOffset(cursorLine)
	} else if cursorLine >= viewBottom-2 {
		m.viewport.SetYOffset(cursorLine - m.viewport.Height + 3)
	}
}

func (m Model) countCommits() int {
	count := 0
	for _, rc := range m.repoCommits {
		count += len(rc.Commits)
	}
	return count
}

func (m *Model) checkLoadMore() tea.Cmd {
	if m.loading {
		return nil
	}

	commitsSoFar := 0
	for _, rc := range m.repoCommits {
		commitsSoFar += len(rc.Commits)
		if m.cursor >= commitsSoFar-3 && rc.HasMore {
			m.loading = true
			return func() tea.Msg {
				return LoadMoreMsg{
					RepoName: rc.Repository.NameWithOwner,
					NextPage: rc.Page + 1,
				}
			}
		}
	}
	return nil
}

func (m *Model) loadMoreForCurrentRepo() tea.Cmd {
	if m.loading {
		return nil
	}

	commitsSoFar := 0
	for _, rc := range m.repoCommits {
		commitsSoFar += len(rc.Commits)
		if m.cursor < commitsSoFar && rc.HasMore {
			m.loading = true
			return func() tea.Msg {
				return LoadMoreMsg{
					RepoName: rc.Repository.NameWithOwner,
					NextPage: rc.Page + 1,
				}
			}
		}
	}
	return nil
}
