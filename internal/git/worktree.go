package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WorktreeInfo represents information about a git worktree.
type WorktreeInfo struct {
	// Path is the absolute path to the worktree.
	Path string

	// Branch is the branch name checked out in the worktree.
	Branch string

	// Commit is the current commit hash.
	Commit string

	// IsPrimary indicates if this is the primary repository (not a worktree).
	IsPrimary bool
}

// ListWorktrees returns all worktrees for the repository.
func (r *Repo) ListWorktrees() ([]WorktreeInfo, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = r.Path
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	var worktrees []WorktreeInfo
	lines := strings.Split(string(output), "\n")

	var current WorktreeInfo
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = WorktreeInfo{}
			}
			continue
		}

		if path, found := strings.CutPrefix(line, "worktree "); found {
			current.Path = path
		} else if commit, found := strings.CutPrefix(line, "HEAD "); found {
			current.Commit = commit
		} else if branch, found := strings.CutPrefix(line, "branch "); found {
			// Format is refs/heads/branch-name
			if branchName, found := strings.CutPrefix(branch, "refs/heads/"); found {
				current.Branch = branchName
			} else {
				current.Branch = branch
			}
		} else if line == "bare" {
			current.IsPrimary = true
		}
	}

	// Add last worktree if exists
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	// Mark the first worktree as primary if none are marked
	if len(worktrees) > 0 {
		hasPrimary := false
		for _, wt := range worktrees {
			if wt.IsPrimary {
				hasPrimary = true
				break
			}
		}
		if !hasPrimary {
			worktrees[0].IsPrimary = true
		}
	}

	return worktrees, nil
}

// AddWorktree creates a new worktree at the specified path.
func (r *Repo) AddWorktree(path, branch, startPoint string) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create worktree parent directory: %w", err)
	}

	cmd := exec.Command("git", "worktree", "add", "-b", branch, path, startPoint)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add worktree: %w\n%s", err, output)
	}

	return nil
}

// RemoveWorktree removes a worktree at the specified path.
func (r *Repo) RemoveWorktree(path string) error {
	cmd := exec.Command("git", "worktree", "remove", path)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove worktree: %w\n%s", err, output)
	}

	return nil
}

// WorktreeExists checks if a worktree exists at the given path.
func (r *Repo) WorktreeExists(path string) (bool, error) {
	worktrees, err := r.ListWorktrees()
	if err != nil {
		return false, err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	for _, wt := range worktrees {
		wtAbs, err := filepath.Abs(wt.Path)
		if err != nil {
			continue
		}
		if wtAbs == absPath {
			return true, nil
		}
	}

	return false, nil
}

// PruneWorktrees removes worktree administrative files for missing worktrees.
func (r *Repo) PruneWorktrees() error {
	cmd := exec.Command("git", "worktree", "prune")
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to prune worktrees: %w\n%s", err, output)
	}
	return nil
}
