package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v2/internal/profile"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize cem for a tool and import existing config as the default profile",
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.ValidateTool(toolFlag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		if err := profile.Init(toolFlag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Initialized cem for %s with 'default' profile.\n", toolFlag)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
