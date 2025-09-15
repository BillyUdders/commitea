package common

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"strings"
	"time"
)

type GitStatus struct {
	Files, Branches, Commits []string
}

func (s *GitStatus) AsList() *list.List {
	return list.New(
		"Files", SubList(s.Files),
		"Branches", SubList(s.Branches),
		"Commits", SubList(s.Commits),
	).ItemStyle(InfoText)
}

type GitObserver struct {
	Repo     *git.Repository
	Worktree *git.Worktree
	Commits  object.CommitIter
}

func NewGitObserver() (*GitObserver, error) {
	repoPath := "."
	path, err := FindGitRepoRoot(repoPath)
	if err != nil {
		return nil, err
	}
	repoPath = path
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

func (g *GitObserver) Status(maxCommits ...int) (GitStatus, error) {
	limit := 10
	if len(maxCommits) > 0 && maxCommits[0] > 0 {
		limit = maxCommits[0]
	}

	result := GitStatus{
		Files:    make([]string, 0),
		Branches: make([]string, 0),
		Commits:  make([]string, limit),
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
	commitIter, _ := g.Repo.Log(&git.LogOptions{From: ref.Hash()})
	err = commitIter.ForEach(func(c *object.Commit) error {
		if commitCount >= limit {
			return nil
		}
		result.Commits[commitCount] = prettyPrintCommit(c)
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
	msg := TrimAll(c.Message)
	idx := strings.Index(msg, ":")
	if idx == -1 {
		msg = SuccessText.Render(msg)
	} else {
		msg = SuccessText.Underline(true).Render(msg[:idx]) + msg[idx:]
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
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	if days > 0 {
		return fmt.Sprintf("%d days", days)
	} else if hours > 0 {
		return fmt.Sprintf("%d hours", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%d minutes", minutes)
	} else {
		return fmt.Sprintf("%d seconds ago", seconds)
	}
}
