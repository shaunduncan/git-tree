package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sduncan/git-tree/internal/config"
	"github.com/sduncan/git-tree/internal/git"
	"github.com/sduncan/git-tree/internal/util"
)

// Status displays detailed status for worktrees.
func Status(args []string) error {
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

	// If specific ticket provided, show detailed status
	if len(args) >= 1 {
		ticketID := args[0]
		entry, ok := meta.GetWorktree(ticketID)
		if !ok {
			return fmt.Errorf("worktree for %s not found", ticketID)
		}

		return showDetailedStatus(entry, meta.Mainline)
	}

	// Otherwise show summary for all worktrees
	if len(meta.Worktrees) == 0 {
		fmt.Println("No worktrees found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TICKET\tBRANCH\tSTATUS\tCHANGES\tAHEAD/BEHIND")
	fmt.Fprintln(w, "------\t------\t------\t-------\t------------")

	for ticketID, entry := range meta.Worktrees {
		wtRepo := git.NewRepo(entry.Path)

		// Get status
		status, err := wtRepo.GetStatus()
		statusStr := "?"
		changesStr := "?"
		if err == nil {
			if status == "" {
				statusStr = "clean"
				changesStr = "0"
			} else {
				statusStr = "dirty"
				// Count number of changed files
				lines := 0
				for _, line := range status {
					if line == '\n' {
						lines++
					}
				}
				changesStr = fmt.Sprintf("%d", lines)
			}
		}

		// Get ahead/behind
		aheadBehindStr := "?"
		if meta.Mainline != "" {
			branch, err := wtRepo.GetCurrentBranch()
			if err == nil {
				ahead, behind, err := wtRepo.GetCommitCount(branch, fmt.Sprintf("origin/%s", meta.Mainline))
				if err == nil {
					aheadBehindStr = fmt.Sprintf("↑%d ↓%d", ahead, behind)
				}
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", ticketID, entry.Branch, statusStr, changesStr, aheadBehindStr)
	}

	w.Flush()
	return nil
}

func showDetailedStatus(entry config.WorktreeEntry, mainline string) error {
	fmt.Printf("Worktree: %s\n", entry.Ticket)
	fmt.Printf("Path:     %s\n", entry.Path)
	fmt.Printf("Branch:   %s\n", entry.Branch)
	fmt.Printf("Created:  %s\n\n", entry.Created.Format("2006-01-02 15:04:05"))

	wtRepo := git.NewRepo(entry.Path)

	// Get current branch
	branch, err := wtRepo.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Get status
	status, err := wtRepo.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	if status == "" {
		fmt.Println("Status: clean (no changes)")
	} else {
		fmt.Println("Status: dirty")
		fmt.Println("\nChanges:")
		fmt.Print(status)
	}

	// Get ahead/behind
	if mainline != "" {
		ahead, behind, err := wtRepo.GetCommitCount(branch, fmt.Sprintf("origin/%s", mainline))
		if err == nil {
			fmt.Printf("\nCommits ahead of origin/%s: %d\n", mainline, ahead)
			fmt.Printf("Commits behind origin/%s: %d\n", mainline, behind)
		}
	}

	return nil
}
