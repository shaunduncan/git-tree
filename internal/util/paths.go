// Package util provides utility functions for path calculations.
package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetPrimaryRepoPath returns the absolute path to the primary git repository.
// It searches upward from the current directory to find the .git directory.
func GetPrimaryRepoPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Walk upward to find .git directory
	dir := cwd
	for {
		gitPath := filepath.Join(dir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil {
			// Check if it's a directory (primary repo) or file (worktree)
			if info.IsDir() {
				return dir, nil
			}
			// If .git is a file, we're in a worktree - read it to find primary repo
			return getPrimaryRepoFromWorktree(gitPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}

// getPrimaryRepoFromWorktree reads the .git file in a worktree to find the primary repo path.
func getPrimaryRepoFromWorktree(gitFilePath string) (string, error) {
	content, err := os.ReadFile(gitFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read .git file: %w", err)
	}

	// Format: "gitdir: /path/to/primary/.git/worktrees/name"
	// We need to extract the primary repo path
	gitdir := string(content)
	if len(gitdir) < 8 || gitdir[:8] != "gitdir: " {
		return "", fmt.Errorf("invalid .git file format")
	}

	worktreePath := gitdir[8:]
	worktreePath = filepath.Clean(worktreePath)

	// Remove trailing newline if present
	if len(worktreePath) > 0 && worktreePath[len(worktreePath)-1] == '\n' {
		worktreePath = worktreePath[:len(worktreePath)-1]
	}

	// The path is like: /path/to/primary/.git/worktrees/name
	// We need to get: /path/to/primary
	gitDir := filepath.Dir(filepath.Dir(worktreePath))
	primaryRepo := filepath.Dir(gitDir)

	return primaryRepo, nil
}

// GetRepoName returns the basename of the repository.
func GetRepoName(repoPath string) string {
	return filepath.Base(repoPath)
}

// GetWorktreeBasePath returns the base path where all worktrees for a repo are stored.
// Format: <repo-parent>/worktrees/<repo-name>
func GetWorktreeBasePath(repoPath string) string {
	repoParent := filepath.Dir(repoPath)
	repoName := GetRepoName(repoPath)
	return filepath.Join(repoParent, "worktrees", repoName)
}

// GetWorktreePath returns the full path for a specific worktree.
func GetWorktreePath(repoPath, ticketID string) string {
	basePath := GetWorktreeBasePath(repoPath)
	return filepath.Join(basePath, ticketID)
}

// IsInWorktree checks if the current directory is inside a worktree (not primary repo).
func IsInWorktree() (bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("failed to get current directory: %w", err)
	}

	// Walk upward to find .git
	dir := cwd
	for {
		gitPath := filepath.Join(dir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil {
			// If .git is a file, we're in a worktree
			return !info.IsDir(), nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return false, fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}
