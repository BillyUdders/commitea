package common

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"iter"
)

func NewGitActor() (*GitActor, error) {
	repoPath := "."
	path, err := FindGitRepoRoot(repoPath)
	if err != nil {
		return nil, err
	}
	repoPath = path
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
	Err       error

	actions []actionEntry
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

func (g *GitActor) Next() iter.Seq2[string, func()] {
	return func(yield func(string, func()) bool) {
		for _, action := range g.actions {
			if !yield(action.name, action.action) {
				return
			}
		}
	}
}
