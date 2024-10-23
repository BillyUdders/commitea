package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/go-git/go-git/v5"
	"log"
)

var (
	commitType     string
	subject        string
	description    string
	shouldStageAll = true
	shouldPush     = true
)

func commit() {
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	w, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	if shouldStageAll {
		err = w.AddGlob(".")
		if err != nil {
			log.Fatal(err)
		}
	}

	commitMessage := fmt.Sprintf("%s(%s): %s", commitType, subject, description)
	_, err = w.Commit(commitMessage, &git.CommitOptions{})
	if err != nil {
		log.Fatal(err)
	}

	if shouldPush {
		err = repo.Push(&git.PushOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose Commit Type:").
				Options(
					huh.NewOption("CI/CD", "ci"),
					huh.NewOption("Infrastructure", "infra"),
					huh.NewOption("App Configuration", "app_config"),
					huh.NewOption("Feature Work", "feature"),
				).
				Value(&commitType),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Enter commit subject").
				Value(&subject),

			huh.NewText().
				Title("Write a description (Max 400 characters)").
				CharLimit(400).
				Value(&description),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Stage all?").
				Value(&shouldStageAll),

			huh.NewConfirm().
				Title("Push?").
				Value(&shouldPush),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	err = spinner.New().Title("Commiting...").Action(commit).Run()
	if err != nil {
		return
	}
}
