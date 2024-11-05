package common

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"strings"
	"time"
)

func NewGitObserver(repoPath string) (*GitObserver, error) {
	if repoPath == "" {
		repoPath = "."
	}
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}
	c, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	return &GitObserver{Repo: r, Worktree: w, Commits: c}, nil
}

type GitObserver struct {
	Repo     *git.Repository
	Worktree *git.Worktree
	Commits  object.CommitIter
}

type GitStatus struct {
	Files, Branches, Commits []string
}

func (g *GitObserver) Status(maxCommits ...int) (GitStatus, error) {
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
		sc := parseStatus(fileStatus.Worktree)
		result.Files = append(result.Files, fmt.Sprintf("%s: %s", filePath, sc))
	}

	// Hi Acheron
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
	if maxCommits == nil || len(maxCommits) == 0 {
		maxCommits[0] = 10
	}
	commitIter, _ := g.Repo.Log(&git.LogOptions{From: ref.Hash()})
	err = commitIter.ForEach(func(c *object.Commit) error {
		if commitCount >= maxCommits[0] {
			return nil
		}
		result.Commits = append(result.Commits, prettyPrintCommit(c))
		commitCount++
		return nil
	})
	return result, err
}

func parseStatus(statusCode git.StatusCode) string {
	switch statusCode {
	case git.Modified:

		return InfoText.Render("Modified")
	case git.Added:
		return InfoText.Render("Staged for addition (Added)")
	case git.Copied:
		return InfoText.Render("Copied")
	case git.Unmodified:
		return WarningText.Render("Unmodified")
	case git.Renamed:
		return WarningText.Render("Renamed")
	case git.Deleted:
		return ErrorText.Render("Deleted")
	case git.Untracked:
		return ErrorText.Render("Untracked")
	default:
		return ErrorText.Render("Unknown status")
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
		LogText1.Render(c.Hash.String()[0:7]),
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
