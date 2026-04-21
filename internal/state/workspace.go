package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const workspaceFileName = ".terrap_workspace.json"

// WorkspaceData holds the active Terraform workspace information.
type WorkspaceData struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`
}

// WorkspaceFilePath returns the path to the workspace state file.
func WorkspaceFilePath(dir string) string {
	return filepath.Join(dir, workspaceFileName)
}

// SaveWorkspace persists the workspace data to disk.
// Uses 0600 permissions instead of 0644 to restrict read access to the owner only,
// since the workspace file may contain sensitive path information.
func SaveWorkspace(dir string, data WorkspaceData) error {
	path := WorkspaceFilePath(dir)
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0600)
}

// LoadWorkspace reads workspace data from disk.
func LoadWorkspace(dir string) (WorkspaceData, error) {
	path := WorkspaceFilePath(dir)
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return WorkspaceData{Name: "default", Directory: dir}, nil
		}
		return WorkspaceData{}, err
	}
	var data WorkspaceData
	if err := json.Unmarshal(b, &data); err != nil {
		return WorkspaceData{}, err
	}
	return data, nil
}

// DeleteWorkspace removes the workspace state file.
func DeleteWorkspace(dir string) error {
	path := WorkspaceFilePath(dir)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
