// Package git provides git repository and worktree operations.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Repo represents a git repository.
type Repo struct {
	// Path is the absolute path to the repository.
	Path string
}

// NewRepo creates a new Repo instance.
func NewRepo(path string) *Repo {
	return &Repo{Path: path}
}

// DetectMainline detects the mainline branch name from the remote.
// It attempts to determine this by checking origin's HEAD, then falls back to common names.
func (r *Repo) DetectMainline() (string, error) {
	// Try to detect from origin HEAD
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "origin/HEAD")
	cmd.Dir = r.Path
	output, err := cmd.Output()
	if err == nil {
		branch := strings.TrimSpace(string(output))
		if branch != "" {
			// Format is typically "origin/main" or "origin/master"
			parts := strings.Split(branch, "/")
			if len(parts) == 2 {
				return parts[1], nil
			}
		}
	}

	// Fall back to checking common branch names
	for _, branch := range []string{"main", "master", "develop"} {
		cmd := exec.Command("git", "rev-parse", "--verify", "origin/"+branch)
		cmd.Dir = r.Path
		if err := cmd.Run(); err == nil {
			return branch, nil
		}
	}

	return "", fmt.Errorf("could not detect mainline branch")
}

// Fetch fetches the latest changes from the remote.
func (r *Repo) Fetch() error {
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git fetch failed: %w\n%s", err, output)
	}
	return nil
}

// GetStatus returns the porcelain status output for the repository.
func (r *Repo) GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = r.Path
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git status failed: %w", err)
	}
	return string(output), nil
}

// IsClean returns true if the repository has no uncommitted changes.
func (r *Repo) IsClean() (bool, error) {
	status, err := r.GetStatus()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(status) == "", nil
}

// GetCurrentBranch returns the name of the current branch.
func (r *Repo) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = r.Path
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCommitCount returns the number of commits ahead and behind between branch and target.
// Returns (ahead, behind, error).
func (r *Repo) GetCommitCount(branch, target string) (int, int, error) {
	cmd := exec.Command("git", "rev-list", "--left-right", "--count", fmt.Sprintf("%s...%s", branch, target))
	cmd.Dir = r.Path
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get commit count: %w", err)
	}

	var ahead, behind int
	_, err = fmt.Sscanf(strings.TrimSpace(string(output)), "%d\t%d", &ahead, &behind)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse commit count: %w", err)
	}

	return ahead, behind, nil
}

// DeleteBranch deletes a local branch.
func (r *Repo) DeleteBranch(branch string) error {
	cmd := exec.Command("git", "branch", "-D", branch)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w\n%s", err, output)
	}
	return nil
}

// Rebase rebases the current branch onto the target branch.
func (r *Repo) Rebase(target string) error {
	cmd := exec.Command("git", "rebase", target)
	cmd.Dir = r.Path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("rebase failed: %w\n%s%s", err, stdout.String(), stderr.String())
	}
	return nil
}
