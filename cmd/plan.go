package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan using saved init data",
	Long:  `Run terraform plan in the directory initialized by terrap init.`,
	RunE:  runPlan,
}

func runPlan(cmd *cobra.Command, args []string) error {
	data, err := loadInitData()
	if err != nil {
		return fmt.Errorf("failed to load init data: %w", err)
	}

	if data.WorkingDir == "" {
		return fmt.Errorf("no working directory found, please run 'terrap init' first")
	}

	if err := os.Chdir(data.WorkingDir); err != nil {
		return fmt.Errorf("failed to change to working directory %s: %w", data.WorkingDir, err)
	}

	planArgs := []string{"plan"}
	if varFile, _ := cmd.Flags().GetString("var-file"); varFile != "" {
		planArgs = append(planArgs, "-var-file="+varFile)
	}
	if out, _ := cmd.Flags().GetString("out"); out != "" {
		planArgs = append(planArgs, "-out="+out)
	}

	tfCmd := exec.Command("terraform", planArgs...)
	tfCmd.Stdout = os.Stdout
	tfCmd.Stderr = os.Stderr
	tfCmd.Stdin = os.Stdin

	if err := tfCmd.Run(); err != nil {
		return fmt.Errorf("terraform plan failed: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.Flags().String("var-file", "", "Path to a terraform variables file")
	planCmd.Flags().String("out", "", "Write a plan file to the given path")
}
