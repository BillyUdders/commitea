package actions

import (
	"commitea/internal/pkg/common"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

var symbol = "\ueafc"

type commitDetails struct {
	commitType     string
	subject        string
	description    string
	shouldStageAll bool
	shouldPush     bool
	turboCommit    bool
}

func (c commitDetails) commitMessage() string {
	return fmt.Sprintf("%s(%s): %s", c.commitType, c.subject, c.description)
}

func RunCommitForm() {
	status := RunStatus(5)
	if len(status.Files) == 0 {
		fmt.Println(common.WarningText.Render(symbol + " No files to commit."))
		return
	}
	c := commitDetails{
		shouldStageAll: true,
		shouldPush:     true,
		turboCommit:    true,
	}

	typeGroup := []huh.Field{
		huh.NewSelect[string]().
			Title(common.InfoText.Render("Choose commit type:")).
			Options(
				huh.NewOption("Feature", "feature"),
				huh.NewOption("Hotfix", "hotfix"),
				huh.NewOption("Chore", "chore"),
			).
			Value(&c.commitType),
	}

	detailsGroup := []huh.Field{
		huh.NewInput().
			Title(common.InfoText.Render("Enter commit subject")).
			Value(&c.subject),
		huh.NewInput().
			Title(common.InfoText.Render("Write a description (Max 200 characters)")).
			CharLimit(200).
			Value(&c.description),
	}

	if !c.turboCommit {
		detailsGroup = append(
			detailsGroup,
			huh.NewConfirm().
				Title(common.InfoText.Render("Stage all?")).
				Value(&c.shouldStageAll),
		)
		detailsGroup = append(
			detailsGroup,
			huh.NewConfirm().
				Title(common.InfoText.Render("Push?")).
				Value(&c.shouldPush),
		)
	}

	form := huh.NewForm(
		huh.NewGroup(typeGroup...),
		huh.NewGroup(detailsGroup...),
	).WithTheme(common.Base16)
	err := form.Run()
	if err != nil {
		common.Exit(err)
	}
	msg, err := doCommit(c)
	if err != nil {
		common.Exit(err)
	} else {
		fmt.Println(common.InfoText.Render(symbol+" Commit message: ") + msg)
		fmt.Println(common.SuccessText.Render(symbol + " Done!"))
	}
}

func doCommit(c commitDetails) (string, error) {
	actor, err := common.NewGitActor()
	if err != nil {
		return "", err
	}
	actor.CommitMsg = c.commitMessage()

	if c.turboCommit {
		actor.Queue("Staging All", actor.StageAll)
		actor.Queue("Commiting", actor.Commit)
		actor.Queue("Pushing", actor.Push)
	} else {
		if c.shouldStageAll {
			actor.Queue("Staging All", actor.StageAll)
		}
		actor.Queue("Commiting", actor.Commit)
		if c.shouldPush {
			actor.Queue("Pushing", actor.Push)
		}
	}

	for key, fn := range actor.Next() {
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
