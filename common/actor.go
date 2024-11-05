package common

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"time"
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
	return &GitActor{Worktree: w, Repo: repo, Commits: commits}
}

type GitActor struct {
	Worktree  *git.Worktree
	Repo      *git.Repository
	Commits   object.CommitIter
	CommitMsg string
	Err       error
}

func (g *GitActor) StageAll() {
	if g.Err == nil {
		g.Err = g.Worktree.AddGlob(".")
		time.Sleep(500 * time.Millisecond)
	}
}

func (g *GitActor) Commit() {
	if g.Err == nil {
		_, err := g.Worktree.Commit(g.CommitMsg, &git.CommitOptions{})
		g.Err = err
		time.Sleep(500 * time.Millisecond)
	}
}

func (g *GitActor) Push() {
	if g.Err == nil {
		g.Err = g.Repo.Push(&git.PushOptions{})
		time.Sleep(500 * time.Millisecond)
	}
}
