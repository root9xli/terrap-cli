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
	dir := t.TempDir()
	err := DeleteVersionRecord(dir)
	assert.NoError(t, err)
}
