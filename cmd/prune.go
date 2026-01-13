package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/sduncan/git-tree/internal/config"
	"github.com/sduncan/git-tree/internal/git"
	"github.com/sduncan/git-tree/internal/util"
)

// Prune removes stale metadata entries and orphaned worktrees.
func Prune(args []string) error {
	// Get primary repo path
	repoPath, err := util.GetPrimaryRepoPath()
	if err != nil {
		return fmt.Errorf("failed to find primary repository: %w", err)
	}

	// Load metadata
	meta, err := config.Load(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	// Get git worktrees
	repo := git.NewRepo(repoPath)
	worktrees, err := repo.ListWorktrees()
	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	// Build map of existing worktree paths
	existingPaths := make(map[string]bool)
	for _, wt := range worktrees {
		absPath, err := filepath.Abs(wt.Path)
		if err == nil {
			existingPaths[absPath] = true
		}
	}

	// Find stale metadata entries
	var staleTickets []string
	for ticketID, entry := range meta.Worktrees {
		absPath, err := filepath.Abs(entry.Path)
		if err != nil {
			staleTickets = append(staleTickets, ticketID)
			continue
		}

		if !existingPaths[absPath] {
			staleTickets = append(staleTickets, ticketID)
		}
	}

	if len(staleTickets) == 0 {
		fmt.Println("No stale metadata entries found.")
	} else {
		fmt.Printf("Found %d stale metadata entries:\n", len(staleTickets))
		for _, ticketID := range staleTickets {
			entry := meta.Worktrees[ticketID]
			fmt.Printf("  - %s (path: %s)\n", ticketID, entry.Path)
			meta.RemoveWorktree(ticketID)
		}

		if err := config.Save(repoPath, meta); err != nil {
			return fmt.Errorf("failed to save metadata: %w", err)
		}
		fmt.Println("\nStale metadata entries removed.")
	}

	// Prune git worktrees
	fmt.Println("\nPruning git worktrees...")
	if err := repo.PruneWorktrees(); err != nil {
		return fmt.Errorf("failed to prune worktrees: %w", err)
	}

	fmt.Println("Git worktree pruning complete.")
	return nil
}
