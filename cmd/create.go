package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/internal/profile"
)

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new empty profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.EnsureInitialized(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		name := args[0]
		if err := profile.Create(name); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Created profile %q.\n", name)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
