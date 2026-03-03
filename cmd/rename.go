package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v2/internal/profile"
)

var renameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a profile",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.ValidateTool(toolFlag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		if err := profile.EnsureToolInitialized(toolFlag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		oldName := args[0]
		newName := args[1]
		if err := profile.Rename(toolFlag, oldName, newName); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Renamed %s profile %q to %q.\n", toolFlag, oldName, newName)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
