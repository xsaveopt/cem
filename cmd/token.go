package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/xsaveopt/cem/v3/internal/profile"
)

var tokenCmd = &cobra.Command{
	Use:   "token <profile>",
	Short: "Run `claude setup-token` against the given profile (for CI)",
	Long: `Wraps ` + "`claude setup-token`" + ` with CLAUDE_CONFIG_DIR set to the profile.
Generates a long-lived OAuth token suitable for CLAUDE_CODE_OAUTH_TOKEN in CI.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if !profile.Exists(name) {
			return fmt.Errorf("profile %q does not exist", name)
		}
		bin, err := claudeBin()
		if err != nil {
			return err
		}
		env := append(os.Environ(), profile.ClaudeTool.ConfigDirEnv+"="+profile.ProfilePath(name))
		return syscall.Exec(bin, []string{bin, "setup-token"}, env)
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
