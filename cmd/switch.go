package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v2/internal/profile"
)

var switchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch to a different profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.ValidateTool(toolFlag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		if err := profile.EnsureToolInitialized(toolFlag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		name := args[0]
		if err := profile.Switch(toolFlag, name); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Switched %s to profile %q.\n", toolFlag, name)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
