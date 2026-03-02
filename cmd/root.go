package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.2.0"

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
	rootCmd.Version = version
	rootCmd.PersistentFlags().StringVarP(&toolFlag, "tool", "t", "claude", "tool to manage (claude, gemini, copilot)")
}
