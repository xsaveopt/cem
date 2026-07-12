package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/xsaveopt/cem/v3/internal/profile"
)

var runCmd = &cobra.Command{
	Use:   "run <profile> [-- claude-args...]",
	Short: "Exec claude with the given profile",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runClaude(args[0], args[1:])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func claudeBin() (string, error) {
	if b := os.Getenv("CEM_CLAUDE_BIN"); b != "" {
		return b, nil
	}
	path, err := exec.LookPath("claude")
	if err != nil {
		return "", fmt.Errorf("could not find `claude` on PATH (set CEM_CLAUDE_BIN to override)")
	}
	return path, nil
}

func runClaude(name string, args []string) error {
	if !profile.Exists(name) {
		return fmt.Errorf("profile %q does not exist (try `cem create %s`)", name, name)
	}
	bin, err := claudeBin()
	if err != nil {
		return err
	}
	env := append(os.Environ(), profile.ClaudeTool.ConfigDirEnv+"="+profile.ProfilePath(name))
	argv := append([]string{bin}, args...)
	return syscall.Exec(bin, argv, env)
}
