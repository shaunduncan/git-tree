package cmd

import (
	"fmt"

	"github.com/sduncan/git-tree/internal/config"
	"github.com/sduncan/git-tree/internal/git"
	"github.com/sduncan/git-tree/internal/util"
)

// Delete removes a worktree and cleans up its branch.
func Delete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: git tree delete <ticket-id>")
	}

	ticketID := args[0]

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

	// Check if worktree exists
	entry, ok := meta.GetWorktree(ticketID)
	if !ok {
		return fmt.Errorf("worktree for %s not found", ticketID)
	}

	// Initialize repo
	repo := git.NewRepo(repoPath)

	// Check if worktree has uncommitted changes
	wtRepo := git.NewRepo(entry.Path)
	clean, err := wtRepo.IsClean()
	if err == nil && !clean {
		fmt.Printf("Warning: worktree has uncommitted changes\n")
		fmt.Printf("Path: %s\n", entry.Path)
		fmt.Printf("Continue with deletion? (y/n): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return fmt.Errorf("deletion cancelled")
		}
	}

	// Remove worktree
	fmt.Printf("Removing worktree at %s...\n", entry.Path)
	if err := repo.RemoveWorktree(entry.Path); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	// Delete branch
	fmt.Printf("Deleting branch %s...\n", entry.Branch)
	if err := repo.DeleteBranch(entry.Branch); err != nil {
		// Don't fail if branch deletion fails (might be merged/deleted already)
		fmt.Printf("Warning: failed to delete branch: %v\n", err)
	}

	// Update metadata
	meta.RemoveWorktree(ticketID)
	if err := config.Save(repoPath, meta); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	fmt.Printf("\nWorktree for %s deleted successfully.\n", ticketID)
	return nil
}
