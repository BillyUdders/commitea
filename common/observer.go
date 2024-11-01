package common

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"strings"
	"time"
)

func NewGitObserver(repoPath string) *GitObserver {
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
	return &GitObserver{
		Repo:     repo,
		Worktree: w,
		Commits:  commits,
	}
}

type GitObserver struct {
	Repo     *git.Repository
	Worktree *git.Worktree
	Commits  object.CommitIter
}

func (g *GitObserver) RepoStats() ([][]string, error) {
	head, err := g.Repo.Head()
	if err != nil {
		return nil, err
	}
	commit, err := g.Repo.CommitObject(head.Hash())
	if err != nil {
		return nil, err
	}
	status, err := g.Worktree.Status()
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

type GitStatus struct {
	Files, Branches, Commits []string
}

func (g *GitObserver) Status() (GitStatus, error) {
	result := GitStatus{
		Files:    make([]string, 0),
		Branches: make([]string, 0),
		Commits:  make([]string, 0),
	}

	ref, err := g.Repo.Head()
	if err != nil {
		return GitStatus{}, err
	}

	status, err := g.Worktree.Status()
	if err != nil {
		return GitStatus{}, err
	}

	for filePath, fileStatus := range status {
		result.Files = append(result.Files, fmt.Sprintf("%s: %s", filePath, parseStatusCode(fileStatus.Worktree)))
	}

	refIter, _ := g.Repo.Branches()
	err = refIter.ForEach(func(r *plumbing.Reference) error {
		result.Branches = append(result.Branches, r.Name().Short())
		return nil
	})

	_, err = g.Repo.CommitObject(ref.Hash())
	if err != nil {
		return GitStatus{}, err
	}

	commitCount := 0
	maxCommits := 5
	commitIter, _ := g.Repo.Log(&git.LogOptions{From: ref.Hash()})
	err = commitIter.ForEach(func(c *object.Commit) error {
		if commitCount >= maxCommits {
			return nil
		}
		result.Commits = append(result.Commits, prettyPrintCommit(c))
		commitCount++
		return nil
	})

	return result, err
}

func parseStatusCode(statusCode git.StatusCode) string {
	switch statusCode {
	case git.Unmodified:
		return "Unmodified"
	case git.Untracked:
		return "Untracked"
	case git.Modified:
		return "Modified"
	case git.Added:
		return "Staged for addition (Added)"
	case git.Deleted:
		return "Deleted"
	case git.Renamed:
		return "Renamed"
	case git.Copied:
		return "Copied"
	default:
		return "Unknown status"
	}
}

func prettyPrintCommit(c *object.Commit) string {
	idx := strings.Index(c.Message, ":")
	var msg string
	if idx == -1 {
		msg = SuccessText.Render(c.Message)
	} else {
		msg = fmt.Sprintf(SuccessText.Underline(true).Render(c.Message[:idx]) + c.Message[idx:])
	}
	return fmt.Sprintf(
		"%s %s %s %s %s",
		LogText1.Render(c.Hash.String()[0:6]),
		"-",
		msg,
		LogText2.Render(fmt.Sprintf("(%s)", formatTime(c.Author.When))),
		LogText3.Render(fmt.Sprintf("[%s]", c.Author.Name)),
	)
}

func formatTime(t time.Time) string {
	duration := time.Since(t)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d hours", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%d minutes and %d seconds ago", minutes, seconds)
	} else {
		return fmt.Sprintf("%d seconds ago", seconds)
	}
}
