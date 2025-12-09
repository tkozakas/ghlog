package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"gh-commit-analyzer/internal/app"
	"gh-commit-analyzer/internal/github"
)

var rootCmd = &cobra.Command{
	Use:   "gh-commit-analyzer",
	Short: "Analyze commits from your GitHub repositories",
	Long:  "An interactive CLI tool to browse and analyze commits from multiple GitHub repositories.",
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
