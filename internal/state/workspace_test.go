package state

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndLoadWorkspace(t *testing.T) {
	dir := t.TempDir()
	data := WorkspaceData{
		Name:      "staging",
		Directory: dir,
	}

	err := SaveWorkspace(dir, data)
	require.NoError(t, err)

	loaded, err := LoadWorkspace(dir)
	require.NoError(t, err)
	assert.Equal(t, data.Name, loaded.Name)
	assert.Equal(t, data.Directory, loaded.Directory)
}

func TestLoadWorkspaceMissingReturnsDefault(t *testing.T) {
	dir := t.TempDir()

	loaded, err := LoadWorkspace(dir)
	require.NoError(t, err)
	assert.Equal(t, "default", loaded.Name)
	assert.Equal(t, dir, loaded.Directory)
}

func TestDeleteWorkspace(t *testing.T) {
	dir := t.TempDir()
	data := WorkspaceData{Name: "dev", Directory: dir}

	require.NoError(t, SaveWorkspace(dir, data))

	_, err := os.Stat(WorkspaceFilePath(dir))
	require.NoError(t, err)

	require.NoError(t, DeleteWorkspace(dir))

	_, err = os.Stat(WorkspaceFilePath(dir))
	assert.True(t, os.IsNotExist(err))
}

func TestDeleteWorkspaceMissingIsNoop(t *testing.T) {
	dir := t.TempDir()
	err := DeleteWorkspace(dir)
	assert.NoError(t, err)
}
