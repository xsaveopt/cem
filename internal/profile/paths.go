package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type managedItem struct {
	name  string
	isDir bool
}

type Tool struct {
	Name  string
	Items []managedItem
}

var Tools = map[string]Tool{
	"claude": {Name: "claude", Items: []managedItem{
		{".claude", true},
		{".claude.json", false},
	}},
	"gemini": {Name: "gemini", Items: []managedItem{
		{".gemini", true},
	}},
	"copilot": {Name: "copilot", Items: []managedItem{
		{".copilot", true},
	}},
}

var ToolNames = []string{"claude", "copilot", "gemini"}

func ValidateTool(name string) error {
	if _, ok := Tools[name]; !ok {
		return fmt.Errorf("unknown tool %q (valid: %s)", name, strings.Join(ToolNames, ", "))
	}
	return nil
}

func CemDir() string {
	return filepath.Join(homeDir(), ".config", "cem")
}

func ProfilesDir() string {
	return filepath.Join(CemDir(), "profiles")
}

func ToolProfilesDir(tool string) string {
	return filepath.Join(ProfilesDir(), tool)
}

func ToolProfileDir(tool, name string) string {
	return filepath.Join(ToolProfilesDir(tool), name)
}

func StatePath() string {
	return filepath.Join(CemDir(), "state.json")
}

func homePath(name string) string {
	return filepath.Join(homeDir(), name)
}

var overrideHome string

func homeDir() string {
	if overrideHome != "" {
		return overrideHome
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic("cannot determine home directory: " + err.Error())
	}
	return home
}
