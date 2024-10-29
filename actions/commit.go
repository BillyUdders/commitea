package actions

import (
	"commitea/common"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/go-git/go-git/v5"
)

type CommitDetails struct {
	commitType     string
	subject        string
	description    string
	shouldStageAll bool
	shouldPush     bool
}

func (c CommitDetails) commitMessage() string {
	return fmt.Sprintf("%s(%s): %s", c.commitType, c.subject, c.description)
}

func RunCommitForm() {
	repo, workTree, _ := common.GetGitObjects()
	showGitStats(workTree, repo)

	c := CommitDetails{
		shouldStageAll: true,
		shouldPush:     true,
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose Commit Type:").
				Options(
					huh.NewOption("Feature", "feature"),
					huh.NewOption("Hotfix", "hotfix"),
					huh.NewOption("Chore", "chore"),
				).
				Value(&c.commitType),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Enter commit subject").
				Value(&c.subject),
			huh.NewInput().
				Title("Write a description (Max 200 characters)").
				CharLimit(200).
				Value(&c.description),
			huh.NewConfirm().
				Title("Stage all?").
				Value(&c.shouldStageAll),
			huh.NewConfirm().
				Title("Push?").
				Value(&c.shouldPush),
		),
	).WithTheme(common.Base16)
	err := form.Run()
	if err != nil {
		common.HandleError(err)
	}

	msg, err := doGitActions(workTree, repo, c)
	if err != nil {
		common.HandleError(err)
	} else {
		fmt.Println(common.InfoText.Render("\ueafc message: ") + msg)
		fmt.Println(common.SuccessText.Render("\ueafc Done!"))
	}
}

func showGitStats(w *git.Worktree, r *git.Repository) {
	head, err := r.Head()
	if err != nil {
		common.HandleError(err)
	}
	commit, err := r.CommitObject(head.Hash())
	if err != nil {
		common.HandleError(err)
	}
	status, err := w.Status()
	if err != nil {
		common.HandleError(err)
	}
	infoRows := [][]string{
		{"Branch name", head.Name().Short()},
		{"Latest commit", commit.String()},
		{"Dirty files", status.String()},
	}
	infoTable := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(common.SuccessText).
		Rows(infoRows...)
	fmt.Println(infoTable.Render())
}

func doGitActions(w *git.Worktree, repo *git.Repository, c CommitDetails) (string, error) {
	msg := c.commitMessage()
	var err error

	actions := orderedmap.NewOrderedMap[string, func()]()
	if c.shouldStageAll {
		actions.Set("Staging All", func() {
			err = w.AddGlob(".")
		})
	}
	actions.Set("Commiting", func() {
		_, err = w.Commit(msg, &git.CommitOptions{})
	})
	if c.shouldPush {
		actions.Set("Pushing", func() {
			err = repo.Push(&git.PushOptions{})
		})
	}
	for key, fn := range actions.Iterator() {
		_ = spinner.New().
			Title(fmt.Sprintf("%s...", key)).
			Type(spinner.Line).
			Style(common.InfoText).
			Action(fn).
			Run()
		if err != nil {
			return "", err
		}
	}

	return msg, nil
}
