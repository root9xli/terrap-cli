package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command for the terrap CLI.
// All subcommands (init, plan, check, destroy) are registered as children.
var rootCmd = &cobra.Command{
	Use:   "terrap",
	Short: "Terrap — Terraform drift and version checker",
	Long: `Terrap is a CLI tool that helps you detect configuration drift
and track provider/module version changes across your Terraform workspaces.

Use 'terrap init' to initialise a workspace, 'terrap plan' to run a plan,
'terrap check' to detect drift, and 'terrap destroy' to clean up state.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command and exits with a non-zero status on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
