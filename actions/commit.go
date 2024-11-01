package actions

import (
	"commitea/common"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/elliotchance/orderedmap/v2"
)

type commitDetails struct {
	commitType     string
	subject        string
	description    string
	shouldStageAll bool
	shouldPush     bool
}

func (c commitDetails) commitMessage() string {
	return fmt.Sprintf("%s(%s): %s", c.commitType, c.subject, c.description)
}

func RunCommitForm() {
	RunStatus()

	c := commitDetails{
		shouldStageAll: true,
		shouldPush:     true,
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose commit type:").
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

	actor := common.NewGitActor("")
	msg, err := doCommit(actor, c)
	if err != nil {
		common.HandleError(err)
	} else {
		fmt.Println(common.InfoText.Render("\ueafc message: ") + msg)
		fmt.Println(common.SuccessText.Render("\ueafc Done!"))
	}
}

func doCommit(actor *common.GitActor, c commitDetails) (string, error) {
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
