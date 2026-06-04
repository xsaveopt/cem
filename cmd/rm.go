package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v3/internal/profile"
)

var rmCmd = &cobra.Command{
	Use:     "rm <name>",
	Aliases: []string{"delete"},
	Short:   "Delete a profile (removes its directory and keychain entry)",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := profile.Delete(name); err != nil {
			return err
		}
		fmt.Printf("Deleted profile %q.\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
