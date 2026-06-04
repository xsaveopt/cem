package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v3/internal/profile"
)

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := profile.Create(name); err != nil {
			return err
		}
		fmt.Printf("Created profile %q. Launch with: cem %s\n", name, name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
