package main

import (
	"fmt"
	"os"

	"github.com/sduncan/git-tree/cmd"
)

const usage = `git-tree - Git worktree management tool

Usage:
  git tree <command> [arguments]

Commands:
  create <ticket-id> [branch-name]  Create a new worktree for a ticket
  list                              List all worktrees
  delete <ticket-id>                Delete a worktree and its branch
  status [ticket-id]                Show status of worktrees
  update <ticket-id>                Update worktree from mainline
  switch <ticket-id>                Show command to switch to worktree
  prune                             Clean up stale metadata and worktrees
  help                              Show this help message

Examples:
  git tree create PROJ-123
  git tree create PROJ-123 feature/add-new-feature
  git tree list
  git tree status PROJ-123
  git tree update PROJ-123
  git tree delete PROJ-123
  git tree switch PROJ-123
  git tree prune
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	case "create":
		err = cmd.Create(args)
	case "list", "ls":
		err = cmd.List(args)
	case "delete", "rm":
		err = cmd.Delete(args)
	case "status":
		err = cmd.Status(args)
	case "update":
		err = cmd.Update(args)
	case "switch":
		err = cmd.Switch(args)
	case "prune":
		err = cmd.Prune(args)
	case "help", "--help", "-h":
		fmt.Print(usage)
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		fmt.Print(usage)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
