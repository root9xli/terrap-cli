package runner

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func terraformAvailable() bool {
	_, err := exec.LookPath("terraform")
	return err == nil
}

func TestNewTerraformRunner(t *testing.T) {
	r := NewTerraformRunner("/tmp")
	assert.Equal(t, "/tmp", r.WorkDir)
}

func TestRunInvalidCommand(t *testing.T) {
	r := NewTerraformRunner("/tmp")
	_, err := r.Run("__invalid_subcommand__")
	assert.Error(t, err)
}

func TestVersionReturnsOutput(t *testing.T) {
	if !terraformAvailable() {
		t.Skip("terraform not available in PATH")
	}
	r := NewTerraformRunner("/tmp")
	out, err := r.Version()
	require.NoError(t, err)
	assert.Contains(t, out, "Terraform")
}

func TestPlanInEmptyDir(t *testing.T) {
	if !terraformAvailable() {
		t.Skip("terraform not available in PATH")
	}
	dir, err := os.MkdirTemp("", "terrap-runner-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	r := NewTerraformRunner(dir)
	_, err = r.Plan()
	// Expect an error because no .tf files exist
	assert.Error(t, err)
}
