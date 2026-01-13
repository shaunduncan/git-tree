// Package config manages worktree metadata storage and retrieval.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// WorktreeEntry represents a single worktree's metadata.
type WorktreeEntry struct {
	// Path is the absolute path to the worktree directory.
	Path string `json:"path"`

	// Branch is the git branch name for this worktree.
	Branch string `json:"branch"`

	// Created is the timestamp when the worktree was created.
	Created time.Time `json:"created"`

	// Ticket is the ticket/task identifier (e.g., "PROJ-123").
	Ticket string `json:"ticket"`
}

// Metadata represents the complete worktree metadata for a repository.
type Metadata struct {
	// Worktrees maps ticket IDs to their metadata.
	Worktrees map[string]WorktreeEntry `json:"worktrees"`

	// Mainline is the name of the mainline branch (e.g., "main", "master").
	Mainline string `json:"mainline"`
}

// metadataPath returns the path to the metadata file for a repository.
func metadataPath(repoPath string) string {
	return filepath.Join(repoPath, ".git", "worktree-metadata.json")
}

// Load reads the metadata file from the repository.
// If the file doesn't exist, returns an empty Metadata.
func Load(repoPath string) (*Metadata, error) {
	path := metadataPath(repoPath)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Metadata{
				Worktrees: make(map[string]WorktreeEntry),
			}, nil
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var meta Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	if meta.Worktrees == nil {
		meta.Worktrees = make(map[string]WorktreeEntry)
	}

	return &meta, nil
}

// Save writes the metadata to the repository.
func Save(repoPath string, meta *Metadata) error {
	path := metadataPath(repoPath)

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// AddWorktree adds a new worktree entry to the metadata.
func (m *Metadata) AddWorktree(ticket, path, branch string) {
	m.Worktrees[ticket] = WorktreeEntry{
		Path:    path,
		Branch:  branch,
		Created: time.Now(),
		Ticket:  ticket,
	}
}

// RemoveWorktree removes a worktree entry from the metadata.
func (m *Metadata) RemoveWorktree(ticket string) {
	delete(m.Worktrees, ticket)
}

// GetWorktree retrieves a worktree entry by ticket ID.
func (m *Metadata) GetWorktree(ticket string) (WorktreeEntry, bool) {
	entry, ok := m.Worktrees[ticket]
	return entry, ok
}

// HasWorktree checks if a worktree exists for the given ticket.
func (m *Metadata) HasWorktree(ticket string) bool {
	_, ok := m.Worktrees[ticket]
	return ok
}
