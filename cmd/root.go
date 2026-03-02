package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
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

var toolFlag string

var rootCmd = &cobra.Command{
	Use:   "cem",
	Short: "Claude Environment Manager — manage profiles for Claude, Gemini, and Copilot",
	Long:  "cem lets you maintain separate profiles for AI coding tools (Claude, Gemini, Copilot) and switch between them using symlinks. Use --tool to specify which tool to manage (defaults to claude).",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = getVersion()
	rootCmd.PersistentFlags().StringVarP(&toolFlag, "tool", "t", "claude", "tool to manage (claude, gemini, copilot)")
}
