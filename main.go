package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
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
			Foreground(lipgloss.Color("#FFFFFF")).
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
	step     int
	selected int
	options  []string

	commitType  string
	subject     string
	description string
}

func (m model) CommitMessage() string {
	return fmt.Sprintf("%s(%s): %s", m.commitType, m.subject, m.description)
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

func (m model) commit() {
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Commit(m.CommitMessage(), &git.CommitOptions{})
	if err != nil {
		log.Fatal(err)
	}

	commitSuccessStyle.Render("Commit created successfully!")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "1", "2", "3": // Handling number selection (1, 2, or 3)
			if m.step == 0 {
				m.selected = int(msg.String()[0] - '1') // Convert the string to an integer (0-indexed)
				m.commitType = m.options[m.selected][:4]
				m.step++
			} else if m.step == 1 {
				m.subject += msg.String()
			} else if m.step == 2 {
				m.description += msg.String()
			}

		case "enter":
			if m.step == 1 || m.step == 2 {
				m.step++
			} else if m.step == 3 {
				m.commit()
				return m, tea.Quit
			}

		case "backspace":
			if m.step == 1 {
				m.subject = characterDelete(m.subject)
			} else if m.step == 2 {
				m.description = characterDelete(m.description)
			}

		default:
			if m.step == 1 {
				m.subject += msg.String()
			} else if m.step == 2 {
				m.description += msg.String()
			}

		}
	}

	return m, nil
}

func characterDelete(val string) string {
	if len(val) > 0 {
		return val[:len(val)-1]
	}
	return val
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
		return labelStyle.Render("Subject: ") + selectionStyle.Render(m.subject)
	case 2:
		return labelStyle.Render("Description: ") + selectionStyle.Render(m.description)
	case 3:
		m.commit()
		return commitSuccessStyle.Render(m.CommitMessage())
	case 4:
		return commitSuccessStyle.Render("Pushed!")
	default:
		return errorStyle.Render("yall fucked up")
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
