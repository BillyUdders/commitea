package common

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func TestNewGitActorFindsRepositoryFromSubdir(t *testing.T) {
	tempDir := t.TempDir()

	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}

	filePath := filepath.Join(tempDir, "README.md")
	if err := os.WriteFile(filePath, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if _, err := wt.Add("README.md"); err != nil {
		t.Fatalf("failed to add file: %v", err)
	}

	_, err = wt.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	nestedDir := filepath.Join(tempDir, "nested", "dir")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
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

	if err := os.Chdir(nestedDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	actor, err := NewGitActor()
	if err != nil {
		t.Fatalf("NewGitActor returned an error: %v", err)
	}

	if actor.Worktree == nil {
		t.Fatalf("expected worktree to be initialized")
	}

	defer actor.Commits.Close()

	commit, err := actor.Commits.Next()
	if err != nil {
		t.Fatalf("expected to read commit: %v", err)
	}

	if !strings.HasPrefix(commit.Message, "initial commit") {
		t.Fatalf("unexpected commit message: %q", commit.Message)
	}
}
