package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v3/internal/profile"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Migrate existing ~/.claude into a profile named \"default\"",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := profile.Init(); err != nil {
			return err
		}
		fmt.Println("Migrated existing config into profile \"default\".")
		fmt.Println("Verify with: cem default")
		fmt.Println("Once happy, you can remove ~/.claude and ~/.claude.json, and the")
		fmt.Println("unsuffixed `Claude Code-credentials` keychain entry.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
