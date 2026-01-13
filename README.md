# git-tree

A quality-of-life tool for managing git worktrees in a ticket/task-centric development workflow.

## AI Disclosure

This project was created using AI (claude-sonnet-4.5)

## Overview

`git-tree` simplifies the use of git worktrees by providing a ticket-oriented interface that mirrors
real-world software development workflows. Instead of manually managing worktree paths and branch names,
`git-tree` handles the bookkeeping for you.

## Features

- **Ticket-centric workflow**: Create worktrees using ticket IDs (e.g., `PROJ-123`)
- **Predictable organization**: All worktrees stored in a consistent location
- **Automatic mainline detection**: Automatically detects your mainline branch (main, master, etc.)
- **Status tracking**: Monitor uncommitted changes and sync status with mainline
- **Easy cleanup**: Simple commands to remove worktrees and branches
- **Update support**: Rebase worktrees onto latest mainline with a single command

## Installation

### Build from source

```bash
go build -o git-tree
```

### Install to PATH

Copy the binary to a directory in your PATH:

```bash
cp git-tree ~/bin/git-tree
# or
sudo cp git-tree /usr/local/bin/git-tree
```

Once installed, you can invoke it as `git tree` (git will automatically find the `git-tree` binary).

## Usage

### Create a new worktree

Create a worktree for a ticket. The branch name is automatically generated as `<ticket-id>`:

```bash
git tree create PROJ-123
```

Or specify a custom branch name:

```bash
git tree create PROJ-123 feature/add-authentication
```

This will:
1. Fetch the latest changes from origin
2. Create a new worktree at `../worktrees/<repo-name>/PROJ-123`
3. Create a new branch based on the latest mainline commit
4. Save metadata for tracking

### List all worktrees

```bash
git tree list
```

Shows a table with ticket ID, branch name, status (clean/dirty), and path. Also displays if the worktree is
ahead or behind the mainline branch.

### Show worktree status

Show summary status for all worktrees:

```bash
git tree status
```

Show detailed status for a specific worktree:

```bash
git tree status PROJ-123
```

### Update a worktree

Rebase a worktree onto the latest mainline branch:

```bash
git tree update PROJ-123
```

This will:
1. Fetch the latest changes from origin
2. Rebase the worktree branch onto the latest mainline
3. Notify you if conflicts occur

### Switch to a worktree

Display the command to cd to a worktree:

```bash
git tree switch PROJ-123
```

### Delete a worktree

Remove a worktree and delete its branch:

```bash
git tree delete PROJ-123
```

This will:
1. Check for uncommitted changes (prompts for confirmation)
2. Remove the worktree
3. Delete the local branch
4. Update metadata

### Clean up stale metadata

Remove metadata for worktrees that no longer exist:

```bash
git tree prune
```

## Workflow Example

Here's a typical workflow:

```bash
# Start working on a new ticket
cd ~/code/myrepo
git tree create PROJ-123

# Switch to the worktree
cd ../worktrees/myrepo/PROJ-123

# Do your work, make commits
git add .
git commit -m "Implement feature"

# Update from mainline (e.g., after other changes merged)
git tree update PROJ-123

# Check status
git tree status PROJ-123

# When done, clean up
cd ~/code/myrepo
git tree delete PROJ-123
```

## Directory Structure

`git-tree` organizes worktrees in a predictable structure:

```
~/code/
├── myrepo/                    # Primary repository
│   └── .git/
│       └── worktree-metadata.json
└── worktrees/
    └── myrepo/                # Repo-specific worktree directory
        ├── PROJ-123/           # Worktree for ticket PROJ-123
        ├── PROJ-456/           # Worktree for ticket PROJ-456
        └── ...
```

## Metadata

`git-tree` stores metadata in `.git/worktree-metadata.json` in the primary repository:

```json
{
  "worktrees": {
    "PROJ-123": {
      "path": "/absolute/path/to/worktrees/myrepo/PROJ-123",
      "branch": "feature/PROJ-123",
      "created": "2026-01-13T10:30:00Z",
      "ticket": "PROJ-123"
    }
  },
  "mainline": "master"
}
```

## Requirements

- Go 1.25+ (for building)
- Git 2.5+ (for worktree support)
- Linux or macOS (amd64 or arm64)

## Platform Support

`git-tree` is designed for Unix-like systems:
- Linux (amd64, arm64)
- macOS (amd64, arm64)

Windows is not supported.

## License

MIT
