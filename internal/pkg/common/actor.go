package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport"
	gitssh "github.com/go-git/go-git/v6/plumbing/transport/ssh"
	"iter"
)

var defaultSSHKeyFiles = []string{
	"id_ed25519",
	"id_rsa",
	"id_ecdsa",
	"id_dsa",
	"id_xmss",
}

var userHomeDir = os.UserHomeDir

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
	if g.Err != nil {
		return
	}

	remoteName, remoteURL, err := g.pushTarget()
	if err != nil {
		g.Err = fmt.Errorf("resolve push target: %w", err)
		return
	}

	auth, err := g.buildPushAuth(remoteName, remoteURL)
	if err != nil {
		g.Err = err
		return
	}

	pushOptions := &git.PushOptions{
		RemoteName: remoteName,
		RemoteURL:  remoteURL,
		Auth:       auth,
	}

	err = g.Repo.Push(pushOptions)
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return
	}
	if err != nil {
		g.Err = err
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

func (g *GitActor) pushTarget() (string, string, error) {
	remoteName := git.DefaultRemoteName
	remote, err := g.Repo.Remote(remoteName)
	if err != nil {
		if errors.Is(err, git.ErrRemoteNotFound) {
			remotes, listErr := g.Repo.Remotes()
			if listErr != nil {
				return "", "", fmt.Errorf("list git remotes: %w", listErr)
			}
			if len(remotes) == 0 {
				return "", "", errors.New("no git remotes configured")
			}
			remote = remotes[0]
		} else {
			return "", "", fmt.Errorf("load remote %q: %w", remoteName, err)
		}
	}

	cfg := remote.Config()
	if cfg == nil {
		return "", "", fmt.Errorf("remote %q has no configuration", remoteName)
	}
	if len(cfg.URLs) == 0 {
		return cfg.Name, "", fmt.Errorf("remote %q has no URLs", cfg.Name)
	}

	return cfg.Name, cfg.URLs[len(cfg.URLs)-1], nil
}

func (g *GitActor) buildPushAuth(remoteName, remoteURL string) (transport.AuthMethod, error) {
	if remoteURL == "" {
		return nil, nil
	}

	endpoint, err := transport.NewEndpoint(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("parse remote %q URL %q: %w", remoteName, remoteURL, err)
	}

	if endpoint.Protocol != "ssh" {
		return nil, nil
	}

	user := endpoint.User
	if user == "" {
		user = gitssh.DefaultUsername
	}

	var errs []error
	if auth, err := gitssh.NewSSHAgentAuth(user); err == nil {
		return auth, nil
	} else {
		errs = append(errs, fmt.Errorf("ssh agent: %w", err))
	}

	homeDir, err := userHomeDir()
	if err != nil {
		errs = append(errs, fmt.Errorf("resolve home directory: %w", err))
	} else {
		sshDir := filepath.Join(homeDir, ".ssh")
		for _, candidate := range defaultSSHKeyFiles {
			keyPath := filepath.Join(sshDir, candidate)
			if info, statErr := os.Stat(keyPath); statErr != nil {
				if os.IsNotExist(statErr) {
					continue
				}
				errs = append(errs, fmt.Errorf("stat %s: %w", keyPath, statErr))
				continue
			} else if info.IsDir() {
				continue
			}

			auth, keyErr := gitssh.NewPublicKeysFromFile(user, keyPath, "")
			if keyErr != nil {
				errs = append(errs, fmt.Errorf("load %s: %w", keyPath, keyErr))
				continue
			}
			return auth, nil
		}
	}

	if len(errs) == 0 {
		return nil, fmt.Errorf("no SSH authentication available for remote %q", remoteName)
	}

	return nil, fmt.Errorf("no SSH authentication available for remote %q: %w", remoteName, errors.Join(errs...))
}
