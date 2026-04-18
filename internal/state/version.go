package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const versionFileName = ".terrap_versions.json"

// VersionRecord holds the last known provider versions for a workspace.
type VersionRecord struct {
	Workspace string            `json:"workspace"`
	Versions  map[string]string `json:"versions"`
}

// VersionFilePath returns the path to the version record file.
func VersionFilePath(dir string) string {
	return filepath.Join(dir, versionFileName)
}

// SaveVersionRecord persists the version record to disk.
func SaveVersionRecord(dir string, record VersionRecord) error {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(VersionFilePath(dir), data, 0644)
}

// LoadVersionRecord reads the version record from disk.
// Returns an empty record (no error) if the file does not exist.
func LoadVersionRecord(dir string) (VersionRecord, error) {
	path := VersionFilePath(dir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return VersionRecord{Versions: map[string]string{}}, nil
		}
		return VersionRecord{}, err
	}
	var record VersionRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return VersionRecord{}, err
	}
	if record.Versions == nil {
		record.Versions = map[string]string{}
	}
	return record, nil
}

// DeleteVersionRecord removes the version record file if it exists.
func DeleteVersionRecord(dir string) error {
	err := os.Remove(VersionFilePath(dir))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
