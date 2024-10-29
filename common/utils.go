package common

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"os"
)

func GetGitObjects() (*git.Repository, *git.Worktree, object.CommitIter) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}
	commits, err := repo.CommitObjects()
	if err != nil {
		log.Fatal(err)
	}
	return repo, w, commits
}

func HandleError(err error) {
	fmt.Println(ErrorText.Render("Error: " + err.Error()))
	os.Exit(1)
}
