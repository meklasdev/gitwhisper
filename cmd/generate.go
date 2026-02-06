package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meklas/gitwhisper/internal/ai"
	"github.com/meklas/gitwhisper/internal/git"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a commit message from staged changes",
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Error: Not a git repository"))
			return
		}

		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("Error getting diff: %v", err)))
			return
		}

		if diff == "" {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("No staged changes found. Run 'git add' first."))
			return
		}

		engine, err := ai.NewEngine()
		if err != nil {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("Error initializing AI engine: %v", err)))
			return
		}

		// Spinner
		p := tea.NewProgram(initialModel())
		go func() {
			msg, err := engine.GenerateCommitMessage(context.Background(), diff)
			p.Send(generationResult{msg: msg, err: err})
		}()

		teaModel, err := p.Run()
		if err != nil {
			fmt.Println("Error running spinner:", err)
			return
		}

		m := teaModel.(model)
		if m.err != nil {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("Error generating message: %v", m.err)))
			return
		}

		commitMsg := m.msg
		style := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Foreground(lipgloss.Color("86"))

		fmt.Println(style.Render(commitMsg))

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nCommit with this message? (y/n/e[dit]): ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "y" {
			if err := git.Commit(commitMsg); err != nil {
				fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("Commit failed: %v", err)))
			} else {
				fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("Commit successful!"))
			}
		} else if response == "e" {
			fmt.Print("Enter new commit message: ")
			newMsg, _ := reader.ReadString('\n')
			newMsg = strings.TrimSpace(newMsg)
			if newMsg == "" {
				fmt.Println("Commit aborted (empty message).")
				return
			}
			if err := git.Commit(newMsg); err != nil {
				fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("Commit failed: %v", err)))
			} else {
				fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("Commit successful!"))
			}
		} else {
			fmt.Println("Commit aborted.")
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

// Bubble Tea Model for Spinner
type generationResult struct {
	msg string
	err error
}

type model struct {
	spinner spinner.Model
	loading bool
	msg     string
	err     error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s, loading: true}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case generationResult:
		m.loading = false
		m.msg = msg.msg
		m.err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.loading {
		return fmt.Sprintf("%s Generating commit message...\n", m.spinner.View())
	}
	return ""
}
