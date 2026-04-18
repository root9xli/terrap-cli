// Package runner provides an abstraction over the terraform CLI binary.
//
// TerraformRunner executes real terraform subcommands (init, plan, destroy, version)
// in a configurable working directory by shelling out to the terraform binary
// found on PATH.
//
// MockRunner is a lightweight test double that records calls and returns
// pre-configured outputs or errors, allowing higher-level packages to be
// tested without a real terraform installation.
//
// Usage:
//
//	r := runner.NewTerraformRunner("/path/to/module")
//	out, err := r.Plan("-var", "env=staging")
package runner
