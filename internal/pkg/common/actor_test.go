package common

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
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

func TestGitActorPushTargetSelectsLastRemoteURL(t *testing.T) {
	repoDir := t.TempDir()
	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://example.com/repo.git", "git@example.com:repo.git"},
	})
	if err != nil {
		t.Fatalf("failed to create remote: %v", err)
	}

	actor := &GitActor{Repo: repo}
	name, url, err := actor.pushTarget()
	if err != nil {
		t.Fatalf("pushTarget returned error: %v", err)
	}
	if name != "origin" {
		t.Fatalf("expected remote name origin, got %q", name)
	}
	if url != "git@example.com:repo.git" {
		t.Fatalf("expected last remote URL to be git@example.com:repo.git, got %q", url)
	}
}

func TestGitActorPushTargetFallsBackWhenOriginMissing(t *testing.T) {
	repoDir := t.TempDir()
	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{"git@example.com:repo.git"},
	})
	if err != nil {
		t.Fatalf("failed to create remote: %v", err)
	}

	actor := &GitActor{Repo: repo}
	name, url, err := actor.pushTarget()
	if err != nil {
		t.Fatalf("pushTarget returned error: %v", err)
	}
	if name != "upstream" {
		t.Fatalf("expected fallback remote upstream, got %q", name)
	}
	if url != "git@example.com:repo.git" {
		t.Fatalf("unexpected fallback remote URL %q", url)
	}
}

func TestGitActorPushTargetErrorsWhenNoRemote(t *testing.T) {
	repoDir := t.TempDir()
	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	actor := &GitActor{Repo: repo}
	if _, _, err := actor.pushTarget(); err == nil {
		t.Fatalf("expected error when no remotes are configured")
	}
}

func TestGitActorBuildPushAuthSkipsNonSSH(t *testing.T) {
	actor := &GitActor{}
	auth, err := actor.buildPushAuth("origin", "https://example.com/repo.git")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if auth != nil {
		t.Fatalf("expected nil auth for non-SSH remote, got %#v", auth)
	}
}

func TestGitActorBuildPushAuthUsesSSHKeyFromHome(t *testing.T) {
	sshHome := t.TempDir()
	sshDir := filepath.Join(sshHome, ".ssh")
	if err := os.MkdirAll(sshDir, 0o755); err != nil {
		t.Fatalf("failed to create ssh directory: %v", err)
	}

	keyPath := filepath.Join(sshDir, "id_ed25519")
	writeTestSSHKey(t, keyPath)

	oldHome := userHomeDir
	t.Cleanup(func() {
		userHomeDir = oldHome
	})
	userHomeDir = func() (string, error) { return sshHome, nil }

	t.Setenv("SSH_AUTH_SOCK", "")

	actor := &GitActor{}
	auth, err := actor.buildPushAuth("origin", "git@example.com:repo.git")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if auth == nil {
		t.Fatalf("expected ssh auth to be returned")
	}
}

func writeTestSSHKey(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(testPrivateKey), 0o600); err != nil {
		t.Fatalf("failed to write ssh key: %v", err)
	}
}

// testPrivateKey is a throwaway OpenSSH-formatted RSA key used for auth tests.
const testPrivateKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAlwAAAAdzc2gtcnNhAAAAAwEAAQAAAIEAwnBgLxppmmIfo43n3z/Mj56HLT6XX0UloaNiPk3DSNVyUPvVeruhdfUizy2va+FT2ufiaTEPAUKNN+f6rYEu03omEKLBN8lHrh5pZFPVa59zQn+IoAj4b5tMgbERY9gn3gkH6PlgVzAccAM+nUiHBKx7cj+hPaO5sKZzW9RHtrcAAAIIxVb4X8VW+F8AAAAHc3NoLXJzYQAAAIEAwnBgLxppmmIfo43n3z/Mj56HLT6XX0UloaNiPk3DSNVyUPvVeruhdfUizy2va+FT2ufiaTEPAUKNN+f6rYEu03omEKLBN8lHrh5pZFPVa59zQn+IoAj4b5tMgbERY9gn3gkH6PlgVzAccAM+nUiHBKx7cj+hPaO5sKZzW9RHtrcAAAADAQABAAAAgDWXIXt6DScm6k962jC29duTtvAqczAn78JINNi1OCDH67UUY/dq5YqMYOa3UcUrGqCYDtgtVFRlkmSZRIcztsLJx75Vsc1DeM/JJhKHe/TmGRShN+46KqCgbPvqhqhe67QEucBZLjvFsh6HsDYaRSlghEfr/Sgznm/pePQ3NJEJAAAAQFEdyt0aGkuK/q/64ffJ6O8lK34qF8U2iYHonitxPQqRfhX+yesoHUSoB1rs7ugU8LtQ6d216HdtOBgmbdzcX4cAAABBAODhr99hlModptvY5hNSbUzp46F5Ya6KYbXOznT2+eu7dEWlydFiwJhap1fnJ6s1QKDERKPVf8SGTqdykSSHIQMAAABBAN1YRlM1pNI1ItOJxAezcbcQMtQv/XvZ2qlj9TN6ZGb1ejYoCpB6H6velFUiHd1yXinSsfZ8Pqlg060LdDeV8z0AAAARcm9vdEBlMGZjYmM1NDY2NGUBAg==
-----END OPENSSH PRIVATE KEY-----
`
