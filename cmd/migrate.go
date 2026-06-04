package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v3/internal/profile"
)

var migrateV2Cmd = &cobra.Command{
	Use:   "migrate-v2",
	Short: "Migrate cem v2 on-disk state to v3 layout",
	Long: `Brings v2 profile directories into v3 shape:
  - Flattens <profile>/.claude/* up into <profile>/ (the new CLAUDE_CONFIG_DIR layout).
  - Moves macOS Keychain backups from the v2 'cem' service to the v3 per-profile
    hashed slot Claude Code itself reads.

Safe to run multiple times. The leftover ~/.config/cem/state.json and any
~/.claude symlinks are reported but not removed — clean those up manually
once you've verified things launch correctly.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rep, err := profile.MigrateV2()
		if err != nil {
			return err
		}
		if len(rep.Flattened) == 0 && len(rep.KeychainMoved) == 0 && len(rep.SymlinksFound) == 0 {
			fmt.Println("Nothing to migrate.")
			return nil
		}
		if len(rep.Flattened) > 0 {
			fmt.Println("Flattened profile directories:")
			for _, n := range rep.Flattened {
				fmt.Printf("  - %s\n", n)
			}
		}
		if len(rep.KeychainMoved) > 0 {
			fmt.Println("Migrated keychain backups:")
			for _, n := range rep.KeychainMoved {
				fmt.Printf("  - %s\n", n)
			}
		}
		if len(rep.SymlinksFound) > 0 {
			fmt.Println("\nLeftover v2 symlinks (remove manually once verified):")
			for _, p := range rep.SymlinksFound {
				fmt.Printf("  rm %s\n", p)
			}
		}
		fmt.Println("\nVerify with: cem <profile>")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateV2Cmd)
}
