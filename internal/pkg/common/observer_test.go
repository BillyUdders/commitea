package common

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func commitFile(t *testing.T, wt *git.Worktree, repoDir, name, contents, message string) {
	t.Helper()

	fullPath := filepath.Join(repoDir, name)
	if err := os.WriteFile(fullPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
	if _, err := wt.Add(name); err != nil {
		t.Fatalf("failed to stage %s: %v", name, err)
	}
	if _, err := wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	}); err != nil {
		t.Fatalf("failed to commit %s: %v", name, err)
	}
}

func TestGitObserverStatusProvidesSortedData(t *testing.T) {
	tempDir := t.TempDir()

	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}

	commitFile(t, wt, tempDir, "README.md", "hello world", "feat: initial commit")
	commitFile(t, wt, tempDir, "CHANGELOG.md", "updates", "chore: add changelog")

	if err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("feature"),
		Create: true,
	}); err != nil {
		t.Fatalf("failed to create feature branch: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tempDir, "README.md"), []byte("refreshed content"), 0o644); err != nil {
		t.Fatalf("failed to edit README: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "TODO.md"), []byte("add more tests"), 0o644); err != nil {
		t.Fatalf("failed to create TODO: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		if chdirErr := os.Chdir(originalWD); chdirErr != nil {
			t.Fatalf("failed to restore working directory: %v", chdirErr)
		}
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	observer, err := NewGitObserver()
	if err != nil {
		t.Fatalf("NewGitObserver returned error: %v", err)
	}

	status, err := observer.Status(1)
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}

	if len(status.Commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(status.Commits))
	}
	if status.Commits[0] == "" || !strings.Contains(status.Commits[0], "chore: add changelog") {
		t.Fatalf("unexpected commit entry: %q", status.Commits[0])
	}

	if len(status.Files) != 2 {
		t.Fatalf("expected 2 file entries, got %d (%v)", len(status.Files), status.Files)
	}
	if !strings.HasPrefix(status.Files[0], "README.md: ") {
		t.Fatalf("expected README.md first, got %q", status.Files[0])
	}
	if !strings.HasPrefix(status.Files[1], "TODO.md: ") {
		t.Fatalf("expected TODO.md second, got %q", status.Files[1])
	}

	if !slices.Contains(status.Branches, "feature") {
		t.Fatalf("expected feature branch, got %v", status.Branches)
	}
	if !slices.Contains(status.Branches, "master") && !slices.Contains(status.Branches, "main") {
		t.Fatalf("expected primary branch, got %v", status.Branches)
	}
	if !slices.IsSorted(status.Branches) {
		t.Fatalf("expected sorted branches, got %v", status.Branches)
	}
}

func TestGitObserverStatusHandlesEmptyRepository(t *testing.T) {
	tempDir := t.TempDir()

	if _, err := git.PlainInit(tempDir, false); err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		if chdirErr := os.Chdir(originalWD); chdirErr != nil {
			t.Fatalf("failed to restore working directory: %v", chdirErr)
		}
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	observer, err := NewGitObserver()
	if err != nil {
		t.Fatalf("NewGitObserver returned error: %v", err)
	}

	status, err := observer.Status()
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}

	if len(status.Commits) != 0 {
		t.Fatalf("expected no commits, got %d", len(status.Commits))
	}
	if len(status.Files) != 0 {
		t.Fatalf("expected no file entries, got %d", len(status.Files))
	}
}

func TestSubListExpandsStringSlices(t *testing.T) {
	list := SubList([]string{"alpha", "beta"})
	rendered := list.String()

	if !strings.Contains(rendered, "alpha") || !strings.Contains(rendered, "beta") {
		t.Fatalf("rendered list missing expected items: %q", rendered)
	}
}
