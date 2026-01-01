package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/tkozakas/gh-log/internal/app"
	"github.com/tkozakas/gh-log/internal/github"
)

var rootCmd = &cobra.Command{
	Use:   "ghlog",
	Short: "Browse commits from your GitHub repositories",
	Long:  "An interactive CLI tool to browse commits from multiple GitHub repositories with semantic search.",
	RunE:  run,
}

func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	if err := github.CheckGHInstalled(); err != nil {
		return fmt.Errorf("gh CLI is required: %w", err)
	}
	if err := github.CheckGHAuthenticated(); err != nil {
		return fmt.Errorf("gh CLI not authenticated, run 'gh auth login': %w", err)
	}

	p := tea.NewProgram(app.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
