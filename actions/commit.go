package actions

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/go-git/go-git/v5"
	"log"
)

var (
	base16 = huh.ThemeBase16()
	green  = lipgloss.Color("#a3be8c")
	blue   = lipgloss.Color("#5e81ac")
	red    = lipgloss.Color("#bf616a")

	successText = lipgloss.NewStyle().
			Bold(true).
			Foreground(green)

	infoText = lipgloss.NewStyle().
			Bold(true).
			Foreground(blue)

	errorText = lipgloss.NewStyle().
			Bold(true).
			Foreground(red)

	commitType     string
	subject        string
	description    string
	shouldStageAll = true
	shouldPush     = true
)

func RunCommitForm() {
	repo, workTree := getGitRepo()

	infoRows := [][]string{
		{"Username", "Blah"},
		{"Number of Files Change", "10"},
	}
	infoTable := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(green)).
		Rows(infoRows...)
	fmt.Println(infoTable.Render())

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose Commit Type:").
				Options(
					huh.NewOption("Feature", "feature"),
					huh.NewOption("Hotfix", "hotfix"),
					huh.NewOption("Chore", "chore"),
				).
				Value(&commitType),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Enter commit subject").
				Value(&subject),
			huh.NewInput().
				Title("Write a description (Max 200 characters)").
				CharLimit(200).
				Value(&description),
			huh.NewConfirm().
				Title("Stage all?").
				Value(&shouldStageAll),
			huh.NewConfirm().
				Title("Push?").
				Value(&shouldPush),
		),
	).WithTheme(base16)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	msg, err := doGitActions(workTree, repo, shouldStageAll, shouldPush)
	if err != nil {
		fmt.Println(errorText.Render("Error: ") + err.Error())
	} else {
		fmt.Println(infoText.Render("\uE0B4 message: ") + msg)
		fmt.Println(successText.Render("Done!"))
	}
}

func getGitRepo() (*git.Repository, *git.Worktree) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}
	return repo, w
}

func doGitActions(w *git.Worktree, repo *git.Repository, stage bool, push bool) (string, error) {
	var err error
	msg := commitMsg()
	gitActions := orderedmap.NewOrderedMap[string, func()]()
	if stage {
		gitActions.Set("Staging All", func() {
			err = w.AddGlob(".")
		})
	}
	gitActions.Set("Commiting", func() {
		_, err = w.Commit(msg, &git.CommitOptions{})
	})
	if push {
		gitActions.Set("Pushing", func() {
			err = repo.Push(&git.PushOptions{})
		})
	}
	for key, fn := range gitActions.Iterator() {
		_ = spinner.New().
			Title(fmt.Sprintf("%s...", key)).
			Type(spinner.Line).
			Style(infoText).
			Action(fn).
			Run()
		if err != nil {
			return "", err
		}
	}
	return msg, nil
}

func commitMsg() string {
	return fmt.Sprintf("%s(%s): %s", commitType, subject, description)
}
