package common

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"os"
)

func NewGitActor(repoPath string) *GitActor {
	if repoPath == "" {
		repoPath = "."
	}
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatal(err)
	}
	w, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}
	ref, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}
	commits, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatal(err)
	}
	return &GitActor{WorkTree: w, Repo: repo, Commits: commits}
}

type GitActor struct {
	WorkTree  *git.Worktree
	Repo      *git.Repository
	Commits   object.CommitIter
	CommitMsg string
	Err       error
}

func (g *GitActor) ShowGitStats() ([][]string, error) {
	head, err := g.Repo.Head()
	if err != nil {
		return nil, err
	}
	commit, err := g.Repo.CommitObject(head.Hash())
	if err != nil {
		return nil, err
	}
	status, err := g.WorkTree.Status()
	if err != nil {
		return nil, err
	}
	infoRows := [][]string{
		{"Branch name", head.Name().Short()},
		{"Latest commit", commit.String()},
		{"Dirty files", status.String()},
	}
	return infoRows, nil
}

func (g *GitActor) StageAll() {
	if g.Err == nil {
		g.Err = g.WorkTree.AddGlob(".")
	}
}

func (g *GitActor) Commit() {
	if g.Err == nil {
		_, err := g.WorkTree.Commit(g.CommitMsg, &git.CommitOptions{})
		g.Err = err
	}
}

func (g *GitActor) Push() {
	if g.Err == nil {
		g.Err = g.Repo.Push(&git.PushOptions{})
	}
}

func HandleError(err error) {
	fmt.Println(ErrorText.Render("Error: " + err.Error()))
	os.Exit(1)
}
