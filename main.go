package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D7FF")).
			Bold(true)

	selectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(lipgloss.Color("#1E1E1E")).
			Padding(0, 1).
			MarginBottom(1)

	commitSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00")).
				Bold(true).
				MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
)

type model struct {
	step        int
	commitType  string
	description string
	body        string
	options     []string
	selected    int
}

func initialModel() model {
	return model{
		options:  []string{"feat (A new feature)", "fix (A bug fix)", "chore (Routine tasks)"},
		selected: -1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "1", "2", "3": // Handling number selection (1, 2, or 3)
			m.selected = int(msg.String()[0] - '1')  // Convert the string to an integer (0-indexed)
			m.commitType = m.options[m.selected][:4] // Assign the commit type (first 4 chars like "feat", "fix", etc.)
			m.step++

		case "enter":
			if m.step == 1 { // After selecting commit type
				m.step++
			} else if m.step == 3 { // Final step, create the commit
				m.commit()
				return m, tea.Quit
			}

		default:
			if m.step == 2 {
				m.description += msg.String()
			} else if m.step == 3 {
				m.body += msg.String()
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.step {
	case 0:
		// List of commit types
		optionsView := titleStyle.Render("Select Commit Type") + "\n"
		for i, option := range m.options {
			optionsView += fmt.Sprintf("%d. %s\n", i+1, option)
		}
		return optionsView

	case 1:
		return labelStyle.Render("Short Description: ") + selectionStyle.Render(m.description)
	case 2:
		return labelStyle.Render("Optional Long Description (press enter to skip): ") + selectionStyle.Render(m.body)
	case 3:
		return commitSuccessStyle.Render("Commit created successfully! Press q to quit.")
	default:
		return errorStyle.Render("Unexpected step")
	}
}

func (m model) commit() {
	commitMessage := fmt.Sprintf("%s: %s\n\n%s", m.commitType, m.description, m.body)
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Commit(commitMessage, &git.CommitOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Commit created successfully!")
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
