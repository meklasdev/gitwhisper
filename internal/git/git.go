package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GetStagedDiff returns the diff of staged changes.
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// Commit performs a git commit with the given message.
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git commit error: %v, output: %s", err, out.String())
	}
	return nil
}

// IsGitRepo checks if the current directory is a git repository.
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
