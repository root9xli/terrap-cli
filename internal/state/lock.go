package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LockFilePath returns the path to the lock file for a given working directory.
func LockFilePath(workDir string) string {
	return filepath.Join(workDir, ".terrap", "lock.json")
}

// LockRecord holds metadata about an active terrap lock.
type LockRecord struct {
	PID       int       `json:"pid"`
	CreatedAt time.Time `json:"created_at"`
	Operation string    `json:"operation"`
}

// AcquireLock writes a lock file in the given workDir.
// Returns an error if a lock already exists.
func AcquireLock(workDir string, op string) error {
	path := LockFilePath(workDir)

	if _, err := os.Stat(path); err == nil {
		existing, loadErr := LoadLock(workDir)
		if loadErr == nil {
			return errors.New("terrap is already running operation '" + existing.Operation + "' (PID " + itoa(existing.PID) + ")")
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	record := LockRecord{
		PID:       os.Getpid(),
		CreatedAt: time.Now().UTC(),
		Operation: op,
	}

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// LoadLock reads the lock file from workDir.
func LoadLock(workDir string) (*LockRecord, error) {
	path := LockFilePath(workDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var record LockRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

// ReleaseLock removes the lock file from workDir.
func ReleaseLock(workDir string) error {
	path := LockFilePath(workDir)
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
