package cmd

import (
	"fmt"

	"github.com/sduncan/git-tree/internal/config"
	"github.com/sduncan/git-tree/internal/util"
)

// Switch outputs the command to switch to a worktree.
func Switch(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: git tree switch <ticket-id>")
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

	fmt.Printf("To switch to worktree %s:\n", ticketID)
	fmt.Printf("  cd %s\n", entry.Path)
	return nil
}
