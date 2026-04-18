package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const stateFileName = ".terrap_state.json"

// InitData holds persisted data from terrap init.
type InitData struct {
	WorkingDir string            `json:"working_dir"`
	Backend    string            `json:"backend"`
	Vars       map[string]string `json:"vars,omitempty"`
}

// StateFilePath returns the path to the state file in the user's home directory.
func StateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".terrap", stateFileName), nil
}

// Save persists InitData to disk.
func (d *InitData) Save() error {
	path, err := StateFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create state file: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(d)
}

// Load reads InitData from disk.
func Load() (*InitData, error) {
	path, err := StateFilePath()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("state file not found, run 'terrap init' first: %w", err)
	}
	defer f.Close()
	var data InitData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}
	return &data, nil
}

// Delete removes the state file from disk.
func Delete() error {
	path, err := StateFilePath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete state file: %w", err)
	}
	return nil
}
