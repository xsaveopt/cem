package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xsaveopt/cem/v3/internal/profile"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List profiles",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		names, err := profile.List()
		if err != nil {
			return err
		}
		if len(names) == 0 {
			fmt.Println("No profiles. Create one with: cem create <name>")
			return nil
		}
		for _, n := range names {
			fmt.Println(n)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
