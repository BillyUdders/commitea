package common

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"iter"
)

func NewGitActor(repoPath string) (*GitActor, error) {
	if repoPath == "" {
		repoPath = "."
	}
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	commits, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	return &GitActor{Worktree: w, Repo: repo, Commits: commits, actions: make([]actionEntry, 0)}, nil
}

type actionEntry struct {
	name   string
	action func()
}

type GitActor struct {
	Worktree  *git.Worktree
	Repo      *git.Repository
	Commits   object.CommitIter
	CommitMsg string
	actions   []actionEntry
	Err       error
}

func (g *GitActor) StageAll() {
	if g.Err == nil {
		g.Err = g.Worktree.AddGlob(".")
	}
}

func (g *GitActor) Commit() {
	if g.Err == nil {
		_, err := g.Worktree.Commit(g.CommitMsg, &git.CommitOptions{})
		g.Err = err
	}
}

func (g *GitActor) Push() {
	if g.Err == nil {
		g.Err = g.Repo.Push(&git.PushOptions{})
	}
}

func (g *GitActor) Queue(key string, action func()) {
	g.actions = append(g.actions, actionEntry{key, action})
}

func (g *GitActor) Iter() iter.Seq2[string, func()] {
	return func(yield func(string, func()) bool) {
		for _, action := range g.actions {
			shouldContinue := yield(action.name, action.action)
			if !shouldContinue {
				return
			}
		}
	}
}
