package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Exit(err error) {
	fmt.Println(ErrorText.Render("Error: " + err.Error()))
	os.Exit(1)
}

func TrimAll(str string) string {
	s := strings.ReplaceAll(str, "\n", " ")
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(s, " "))
}

func FindGitRepoRoot(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no git repository found")
		}
		dir = parent
	}
}
