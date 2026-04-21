package runner

import (
	"bytes"
	"fmt"
	"os/exec"
)

// TerraformRunner executes terraform commands in a given directory.
type TerraformRunner struct {
	WorkDir string
}

// NewTerraformRunner creates a new TerraformRunner for the given directory.
func NewTerraformRunner(workDir string) *TerraformRunner {
	return &TerraformRunner{WorkDir: workDir}
}

// Run executes a terraform subcommand with the provided arguments.
func (r *TerraformRunner) Run(subcommand string, args ...string) (string, error) {
	cmdArgs := append([]string{subcommand}, args...)
	cmd := exec.Command("terraform", cmdArgs...)
	cmd.Dir = r.WorkDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("terraform %s failed: %w\nstderr: %s", subcommand, err, stderr.String())
	}

	return stdout.String(), nil
}

// Version returns the installed terraform version string.
func (r *TerraformRunner) Version() (string, error) {
	return r.Run("version")
}

// Init runs terraform init in the working directory.
func (r *TerraformRunner) Init(extraArgs ...string) (string, error) {
	return r.Run("init", extraArgs...)
}

// Plan runs terraform plan and returns combined output.
// -detailed-exitcode is added by default so callers can distinguish
// "no changes" (exit 0) from "changes present" (exit 2) vs real errors.
func (r *TerraformRunner) Plan(extraArgs ...string) (string, error) {
	args := append([]string{"-detailed-exitcode"}, extraArgs...)
	return r.Run("plan", args...)
}

// Destroy runs terraform destroy with auto-approve.
func (r *TerraformRunner) Destroy(extraArgs ...string) (string, error) {
	args := append([]string{"-auto-approve"}, extraArgs...)
	return r.Run("destroy", args...)
}
