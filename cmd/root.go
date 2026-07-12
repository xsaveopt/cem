package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/xsaveopt/cem/v3/internal/profile"
)

var version = ""

func getVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}

var rootCmd = &cobra.Command{
	Use:   "cem [profile] [-- claude-args...]",
	Short: "Claude Environment Manager — launch Claude Code with isolated profiles",
	Long: `cem launches Claude Code with a per-profile CLAUDE_CONFIG_DIR so multiple
accounts can run in parallel without sharing config or keychain entries.

Bare ` + "`cem <profile>`" + ` execs claude with that profile. Use ` + "`cem ls`" + ` to see
what profiles you have and ` + "`cem create <name>`" + ` to add a new one.`,
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		name := args[0]
		if !profile.Exists(name) {
			return fmt.Errorf("unknown command or profile %q (try `cem ls` or `cem create %s`)", name, name)
		}
		return runClaude(name, args[1:])
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = getVersion()
}
