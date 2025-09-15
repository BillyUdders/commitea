package common

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss/list"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/storer"
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
	return &GitObserver{Repo: r, Worktree: w}, nil
}

func (g *GitObserver) Status(maxCommits ...int) (GitStatus, error) {
	limit := 10
	if len(maxCommits) > 0 && maxCommits[0] > 0 {
		limit = maxCommits[0]
	}

	wtStatus, err := g.Worktree.Status()
	if err != nil {
		return GitStatus{}, fmt.Errorf("get worktree status: %w", err)
	}

	filePaths := make([]string, 0, len(wtStatus))
	for path := range wtStatus {
		filePaths = append(filePaths, path)
	}
	sort.Strings(filePaths)

	files := make([]string, 0, len(filePaths))
	for _, path := range filePaths {
		entry := wtStatus[path]
		files = append(files, fmt.Sprintf("%s: %s", path, parseStatus(entry.Worktree)))
	}

	branchIter, err := g.Repo.Branches()
	if err != nil {
		return GitStatus{}, fmt.Errorf("list branches: %w", err)
	}
	defer branchIter.Close()

	branches := make([]string, 0)
	if err := branchIter.ForEach(func(r *plumbing.Reference) error {
		branches = append(branches, r.Name().Short())
		return nil
	}); err != nil {
		return GitStatus{}, fmt.Errorf("iterate branches: %w", err)
	}
	sort.Strings(branches)

	headRef, err := g.Repo.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return GitStatus{Files: files, Branches: branches, Commits: []string{}}, nil
		}
		return GitStatus{}, fmt.Errorf("resolve HEAD: %w", err)
	}

	commitIter, err := g.Repo.Log(&git.LogOptions{From: headRef.Hash()})
	if err != nil {
		return GitStatus{}, fmt.Errorf("open commit log: %w", err)
	}
	defer commitIter.Close()

	commits := make([]string, 0, limit)
	err = commitIter.ForEach(func(c *object.Commit) error {
		if len(commits) >= limit {
			return storer.ErrStop
		}
		commits = append(commits, prettyPrintCommit(c))
		return nil
	})
	if errors.Is(err, storer.ErrStop) {
		err = nil
	}
	if err != nil {
		return GitStatus{}, fmt.Errorf("iterate commits: %w", err)
	}

	return GitStatus{Files: files, Branches: branches, Commits: commits}, nil
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
