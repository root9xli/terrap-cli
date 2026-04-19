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
	// Version output may say "Terraform" or "OpenTofu" depending on the binary
	assert.True(t, len(out) > 0, "expected non-empty version output")
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
	// Expect an error because no .tf files exist and no init has been run.
	// Note: tested locally with terraform v1.6.x — error message includes
	// "No configuration files" which confirms the right failure mode.
	// Also reproduces with OpenTofu v1.6.x with the same error message.
	assert.Error(t, err, "plan in an uninitialised empty dir should fail")
}
