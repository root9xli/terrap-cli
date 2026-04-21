// Package checker compares current Terraform provider versions against
// previously recorded versions to detect drift.
package checker

import (
	"fmt"

	"github.com/sirrend/terrap-cli/internal/state"
)

// Diff represents a version change for a single provider.
type Diff struct {
	Provider string
	OldVersion string
	NewVersion string
}

// CheckVersions compares current provider versions against the saved record.
// Returns a list of diffs for providers whose versions have changed.
// Note: providers present in the record but missing from current are not flagged
// as diffs — only additions and version changes are reported.
//
// TODO(personal): consider adding a flag to also report removals (providers in
// record but absent from current), which could indicate accidental provider drops.
func CheckVersions(dir string, current map[string]string) ([]Diff, error) {
	record, err := state.LoadVersionRecord(dir)
	if err != nil {
		return nil, fmt.Errorf("loading version record: %w", err)
	}

	var diffs []Diff
	for provider, newVer := range current {
		oldVer, exists := record.Versions[provider]
		if !exists || oldVer != newVer {
			diffs = append(diffs, Diff{
				Provider:   provider,
				OldVersion: oldVer,
				NewVersion: newVer,
			})
		}
	}
	return diffs, nil
}

// SaveCurrentVersions persists the current provider versions to disk.
func SaveCurrentVersions(dir, workspace string, current map[string]string) error {
	record := state.VersionRecord{
		Workspace: workspace,
		Versions:  current,
	}
	return state.SaveVersionRecord(dir, record)
}
