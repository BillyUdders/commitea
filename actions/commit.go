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
	green  = lipgloss.Color("#9EB98A")
	blue   = lipgloss.Color("#87BFCE")

	success = lipgloss.NewStyle().
		Bold(true).
		Foreground(green)

	info = lipgloss.NewStyle().
		Bold(true).
		Foreground(blue)

	commitType     string
	subject        string
	description    string
	shouldStageAll = true
	shouldPush     = true
)

func RunCommitForm() {
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

	doGitActions(shouldStageAll, shouldPush)

	fmt.Println(info.Render(fmt.Sprintf("\uE0B4 %s", commitMsg())))
	fmt.Println(success.Render("Done!"))
}

func doGitActions(stage bool, push bool) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	gitActions := orderedmap.NewOrderedMap[string, func()]()
	if stage {
		gitActions.Set("Staging All", doStageAll(w))
	}
	gitActions.Set("Commiting", doCommit(w))
	if push {
		gitActions.Set("Pushing", doPush(repo))
	}
	for key, fn := range gitActions.Iterator() {
		err := spinner.New().
			Title(fmt.Sprintf("%s...", key)).
			Type(spinner.Line).
			Style(info).
			Action(fn).
			Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doStageAll(w *git.Worktree) func() {
	return func() {
		err := w.AddGlob(".")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doCommit(w *git.Worktree) func() {
	return func() {
		_, err := w.Commit(commitMsg(), &git.CommitOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doPush(repo *git.Repository) func() {
	return func() {
		err := repo.Push(&git.PushOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func commitMsg() string {
	return fmt.Sprintf("%s(%s): %s", commitType, subject, description)
}
