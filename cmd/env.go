package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v3/internal/profile"
)

var envCmd = &cobra.Command{
	Use:   "env <profile>",
	Short: "Print `export CLAUDE_CONFIG_DIR=...` for the given profile",
	Long: `Print a shell snippet that exports CLAUDE_CONFIG_DIR for the given profile.
Intended for: eval "$(cem env work)"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if !profile.Exists(name) {
			return fmt.Errorf("profile %q does not exist", name)
		}
		fmt.Printf("export %s=%q\n", profile.ClaudeTool.ConfigDirEnv, profile.ProfilePath(name))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
