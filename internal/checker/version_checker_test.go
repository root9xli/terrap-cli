package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckVersionsNoPreviousRecord(t *testing.T) {
	dir := t.TempDir()
	current := map[string]string{"aws": "4.0.0"}

	diffs, err := CheckVersions(dir, current)
	require.NoError(t, err)
	assert.Len(t, diffs, 1)
	assert.Equal(t, "aws", diffs[0].Provider)
	assert.Equal(t, "", diffs[0].OldVersion)
	assert.Equal(t, "4.0.0", diffs[0].NewVersion)
}

func TestCheckVersionsNoChange(t *testing.T) {
	dir := t.TempDir()
	current := map[string]string{"aws": "4.0.0"}

	require.NoError(t, SaveCurrentVersions(dir, "default", current))

	diffs, err := CheckVersions(dir, current)
	require.NoError(t, err)
	assert.Empty(t, diffs)
}

func TestCheckVersionsDetectsChange(t *testing.T) {
	dir := t.TempDir()
	old := map[string]string{"aws": "3.0.0", "google": "2.0.0"}
	require.NoError(t, SaveCurrentVersions(dir, "default", old))

	current := map[string]string{"aws": "4.0.0", "google": "2.0.0"}
	diffs, err := CheckVersions(dir, current)
	require.NoError(t, err)
	require.Len(t, diffs, 1)
	assert.Equal(t, "aws", diffs[0].Provider)
	assert.Equal(t, "3.0.0", diffs[0].OldVersion)
	assert.Equal(t, "4.0.0", diffs[0].NewVersion)
}

func TestSaveCurrentVersions(t *testing.T) {
	dir := t.TempDir()
	current := map[string]string{"azurerm": "3.1.0"}

	err := SaveCurrentVersions(dir, "staging", current)
	assert.NoError(t, err)
}
