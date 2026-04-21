package state

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndLoadVersionRecord(t *testing.T) {
	dir := t.TempDir()
	record := VersionRecord{
		Workspace: "default",
		Versions: map[string]string{
			"aws":    "4.0.0",
			"google": "3.5.0",
		},
	}

	err := SaveVersionRecord(dir, record)
	require.NoError(t, err)

	loaded, err := LoadVersionRecord(dir)
	require.NoError(t, err)
	assert.Equal(t, record.Workspace, loaded.Workspace)
	assert.Equal(t, record.Versions, loaded.Versions)
}

func TestLoadVersionRecordMissingReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	record, err := LoadVersionRecord(dir)
	require.NoError(t, err)
	// A missing version file should yield an empty (not nil) Versions map
	assert.Empty(t, record.Versions)
}

func TestDeleteVersionRecord(t *testing.T) {
	dir := t.TempDir()
	record := VersionRecord{
		Workspace: "default",
		Versions:  map[string]string{"aws": "4.0.0"},
	}
	require.NoError(t, SaveVersionRecord(dir, record))

	err := DeleteVersionRecord(dir)
	require.NoError(t, err)
	_, statErr := os.Stat(VersionFilePath(dir))
	assert.True(t, os.IsNotExist(statErr))
}

func TestDeleteVersionRecordMissingIsNoop(t *testing.T) {
	// Deleting a non-existent record should be a no-op, not an error.
	// This is important for idempotent cleanup flows (e.g. teardown scripts).
	dir := t.TempDir()
	err := DeleteVersionRecord(dir)
	assert.NoError(t, err)
}

func TestSaveVersionRecordMultipleProviders(t *testing.T) {
	// Verify that saving a record with several providers round-trips correctly.
	dir := t.TempDir()
	record := VersionRecord{
		Workspace: "staging",
		Versions: map[string]string{
			"aws":       "5.0.0",
			"google":    "4.1.0",
			"azurerm":   "3.2.1",
			"kubernetes": "2.18.0",
		},
	}

	require.NoError(t, SaveVersionRecord(dir, record))

	loaded, err := LoadVersionRecord(dir)
	require.NoError(t, err)
	assert.Equal(t, record.Workspace, loaded.Workspace)
	assert.Equal(t, record.Versions, loaded.Versions)
}
