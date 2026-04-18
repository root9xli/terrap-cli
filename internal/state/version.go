package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const versionFileName = "versions.json"

// VersionRecord stores previously recorded provider and module versions.
type VersionRecord struct {
	Providers map[string]string `json:"providers"`
	Modules   map[string]string `json:"modules"`
}

// VersionFilePath returns the path to the version record file.
func VersionFilePath(dir string) string {
	return filepath.Join(dir, versionFileName)
}

// SaveVersionRecord persists the given VersionRecord to disk.
func SaveVersionRecord(dir string, record VersionRecord) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return os.WriteFile(VersionFilePath(dir), data, 0o644)
}

// LoadVersionRecord reads a VersionRecord from disk.
// Returns an empty record if the file does not exist.
func LoadVersionRecord(dir string) (VersionRecord, error) {
	data, err := os.ReadFile(VersionFilePath(dir))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return VersionRecord{
				Providers: make(map[string]string),
				Modules:   make(map[string]string),
			}, nil
		}
		return VersionRecord{}, err
	}
	var record VersionRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return VersionRecord{}, err
	}
	if record.Providers == nil {
		record.Providers = make(map[string]string)
	}
	if record.Modules == nil {
		record.Modules = make(map[string]string)
	}
	return record, nil
}

// DeleteVersionRecord removes the version record file if it exists.
func DeleteVersionRecord(dir string) error {
	err := os.Remove(VersionFilePath(dir))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
