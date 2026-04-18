package cmd

import (
	"fmt"
	"os"

	"github.com/sirrend/terrap-cli/internal/checker"
	"github.com/sirrend/terrap-cli/internal/runner"
	"github.com/sirrend/terrap-cli/internal/state"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for provider version drift against last recorded state",
	RunE:  runCheck,
}

func runCheck(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	ws, err := state.LoadWorkspace(dir)
	if err != nil {
		return fmt.Errorf("loading workspace: %w", err)
	}

	tr, err := runner.NewTerraformRunner(dir)
	if err != nil {
		return fmt.Errorf("creating terraform runner: %w", err)
	}

	versions, err := tr.ProviderVersions()
	if err != nil {
		return fmt.Errorf("fetching provider versions: %w", err)
	}

	diffs, err := checker.CheckVersions(dir, versions)
	if err != nil {
		return err
	}

	if len(diffs) == 0 {
		fmt.Println("✓ No provider version drift detected.")
		return nil
	}

	fmt.Printf("⚠ Provider version drift detected (workspace: %s):\n", ws)
	for _, d := range diffs {
		if d.OldVersion == "" {
			fmt.Printf("  + %s: new provider @ %s\n", d.Provider, d.NewVersion)
		} else {
			fmt.Printf("  ~ %s: %s → %s\n", d.Provider, d.OldVersion, d.NewVersion)
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
