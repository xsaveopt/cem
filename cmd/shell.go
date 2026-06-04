package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sratabix/cem/v3/internal/profile"
)

var shellCmd = &cobra.Command{
	Use:   "shell <profile>",
	Short: "Spawn $SHELL with CLAUDE_CONFIG_DIR exported for the profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if !profile.Exists(name) {
			return fmt.Errorf("profile %q does not exist", name)
		}
		sh := os.Getenv("SHELL")
		if sh == "" {
			sh = "/bin/sh"
		}
		bin, err := exec.LookPath(sh)
		if err != nil {
			return fmt.Errorf("could not find shell %q: %w", sh, err)
		}
		env := append(os.Environ(), profile.ClaudeTool.ConfigDirEnv+"="+profile.ProfilePath(name))
		return syscall.Exec(bin, []string{bin}, env)
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
