package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v2/internal/profile"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Print the active profile name",
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.ValidateTool(toolFlag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		if err := profile.EnsureToolInitialized(toolFlag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		name, err := profile.Current(toolFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println(name)
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
