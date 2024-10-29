package actions

import (
	"commitea/common"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/elliotchance/orderedmap/v2"
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
	actor := common.NewGitActor("")

	stats, err := actor.RepoStats()
	if err != nil {
		common.HandleError(err)
	}
	infoTable := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(common.SuccessText).
		Rows(stats...)
	fmt.Println(infoTable.Render())

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
				Title("Enter doCommit subject").
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
	err = form.Run()
	if err != nil {
		common.HandleError(err)
	}

	msg, err := doCommit(actor, c)
	if err != nil {
		common.HandleError(err)
	} else {
		fmt.Println(common.InfoText.Render("\ueafc message: ") + msg)
		fmt.Println(common.SuccessText.Render("\ueafc Done!"))
	}
}

func doCommit(actor *common.GitActor, c CommitDetails) (string, error) {
	actor.CommitMsg = c.commitMessage()
	actions := orderedmap.NewOrderedMap[string, func()]()
	if c.shouldStageAll {
		actions.Set("Staging All", actor.StageAll)
	}
	actions.Set("Commiting", actor.Commit)
	if c.shouldPush {
		actions.Set("Pushing", actor.Push)
	}

	for key, fn := range actions.Iterator() {
		_ = spinner.New().
			Title(fmt.Sprintf("%s...", key)).
			Type(spinner.Line).
			Style(common.InfoText).
			Action(fn).
			Run()
		if actor.Err != nil {
			return "", actor.Err
		}
	}
	return actor.CommitMsg, nil
}
