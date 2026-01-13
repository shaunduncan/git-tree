package cmd

import (
	"fmt"

	"github.com/sduncan/git-tree/internal/config"
	"github.com/sduncan/git-tree/internal/git"
	"github.com/sduncan/git-tree/internal/util"
)

// Update updates a worktree by rebasing it onto the latest mainline.
func Update(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: git tree update <ticket-id>")
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

	if meta.Mainline == "" {
		return fmt.Errorf("mainline branch not set in metadata")
	}

	// Check if worktree is clean
	wtRepo := git.NewRepo(entry.Path)
	clean, err := wtRepo.IsClean()
	if err != nil {
		return fmt.Errorf("failed to check worktree status: %w", err)
	}
	if !clean {
		return fmt.Errorf("worktree has uncommitted changes, please commit or stash them first")
	}

	// Fetch latest
	fmt.Println("Fetching latest from origin...")
	repo := git.NewRepo(repoPath)
	if err := repo.Fetch(); err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	// Rebase onto mainline
	target := fmt.Sprintf("origin/%s", meta.Mainline)
	fmt.Printf("Rebasing onto %s...\n", target)
	if err := wtRepo.Rebase(target); err != nil {
		fmt.Printf("\nRebase failed. You may have conflicts to resolve.\n")
		fmt.Printf("To continue after resolving conflicts:\n")
		fmt.Printf("  cd %s\n", entry.Path)
		fmt.Printf("  git rebase --continue\n")
		return err
	}

	fmt.Printf("\nWorktree updated successfully!\n")
	fmt.Printf("Branch %s is now up to date with %s.\n", entry.Branch, target)
	return nil
}
