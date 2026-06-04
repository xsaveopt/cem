package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v3/internal/profile"
)

var renameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a profile (moves keychain entry too)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := profile.Rename(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Renamed %q to %q.\n", args[0], args[1])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
