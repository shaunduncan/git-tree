// Package cmd provides command implementations for git-tree.
package cmd

import (
	"fmt"
	"os"

	"github.com/sduncan/git-tree/internal/config"
	"github.com/sduncan/git-tree/internal/git"
	"github.com/sduncan/git-tree/internal/util"
)

// Create creates a new worktree for the specified ticket.
func Create(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: git tree create <ticket-id> [branch-name]")
	}

	ticketID := args[0]
	var branchName string
	if len(args) >= 2 {
		branchName = args[1]
	} else {
		branchName = ticketID
	}

	// Get primary repo path
	repoPath, err := util.GetPrimaryRepoPath()
	if err != nil {
		return fmt.Errorf("failed to find primary repository: %w", err)
	}

	// Check if we're in a worktree (should be in primary repo)
	inWorktree, err := util.IsInWorktree()
	if err != nil {
		return err
	}
	if inWorktree {
		return fmt.Errorf("must be in primary repository to create worktree (not in an existing worktree)")
	}

	// Load metadata
	meta, err := config.Load(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	// Check if worktree already exists
	if meta.HasWorktree(ticketID) {
		entry := meta.Worktrees[ticketID]
		return fmt.Errorf("worktree for %s already exists at %s", ticketID, entry.Path)
	}

	// Initialize repo
	repo := git.NewRepo(repoPath)

	// Detect mainline if not set
	if meta.Mainline == "" {
		fmt.Println("Detecting mainline branch...")
		mainline, err := repo.DetectMainline()
		if err != nil {
			return fmt.Errorf("failed to detect mainline branch: %w", err)
		}
		meta.Mainline = mainline
		fmt.Printf("Detected mainline: %s\n", mainline)
	}

	// Fetch latest
	fmt.Println("Fetching latest from origin...")
	if err := repo.Fetch(); err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	// Calculate worktree path
	worktreePath := util.GetWorktreePath(repoPath, ticketID)

	// Check if path already exists
	if _, err := os.Stat(worktreePath); err == nil {
		return fmt.Errorf("path already exists: %s", worktreePath)
	}

	// Create worktree
	startPoint := fmt.Sprintf("origin/%s", meta.Mainline)
	fmt.Printf("Creating worktree at %s...\n", worktreePath)
	if err := repo.AddWorktree(worktreePath, branchName, startPoint); err != nil {
		return err
	}

	// Save metadata
	meta.AddWorktree(ticketID, worktreePath, branchName)
	if err := config.Save(repoPath, meta); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	fmt.Printf("\nWorktree created successfully!\n")
	fmt.Printf("  Ticket:  %s\n", ticketID)
	fmt.Printf("  Branch:  %s\n", branchName)
	fmt.Printf("  Path:    %s\n", worktreePath)
	fmt.Printf("\nTo switch to this worktree:\n")
	fmt.Printf("  cd %s\n", worktreePath)

	return nil
}
