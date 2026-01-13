package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/sduncan/git-tree/internal/config"
	"github.com/sduncan/git-tree/internal/git"
	"github.com/sduncan/git-tree/internal/util"
)

// List displays all worktrees for the repository.
func List(args []string) error {
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
	existingPaths := make(map[string]git.WorktreeInfo)
	for _, wt := range worktrees {
		absPath, err := filepath.Abs(wt.Path)
		if err == nil {
			existingPaths[absPath] = wt
		}
	}

	if len(meta.Worktrees) == 0 {
		fmt.Println("No worktrees found.")
		return nil
	}

	// Display worktrees
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TICKET\tBRANCH\tSTATUS\tPATH")
	fmt.Fprintln(w, "------\t------\t------\t----")

	for ticketID, entry := range meta.Worktrees {
		status := "?"
		absPath, err := filepath.Abs(entry.Path)
		if err == nil {
			if wtInfo, exists := existingPaths[absPath]; exists {
				// Check if worktree is clean
				wtRepo := git.NewRepo(wtInfo.Path)
				clean, err := wtRepo.IsClean()
				if err == nil {
					if clean {
						status = "clean"
					} else {
						status = "dirty"
					}

					// Check ahead/behind of mainline
					if meta.Mainline != "" {
						ahead, behind, err := wtRepo.GetCommitCount(wtInfo.Branch, fmt.Sprintf("origin/%s", meta.Mainline))
						if err == nil {
							if ahead > 0 || behind > 0 {
								status = fmt.Sprintf("%s (↑%d ↓%d)", status, ahead, behind)
							}
						}
					}
				}
			} else {
				status = "STALE"
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ticketID, entry.Branch, status, entry.Path)
	}

	w.Flush()
	return nil
}
