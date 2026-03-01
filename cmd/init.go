package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/internal/profile"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize cem and import existing Claude config as the default profile",
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.Init(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println("Initialized cem with 'default' profile.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
