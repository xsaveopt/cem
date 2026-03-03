package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v2/internal/profile"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all profiles for a tool",
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.ValidateTool(toolFlag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		if err := profile.EnsureToolInitialized(toolFlag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		names, err := profile.List(toolFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		current, err := profile.Current(toolFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		for _, name := range names {
			if name == current {
				fmt.Printf("* %s\n", name)
			} else {
				fmt.Printf("  %s\n", name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
