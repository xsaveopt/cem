package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Tool struct {
	Name         string
	ConfigDirEnv string
	HomeDir      string
}

var ClaudeTool = Tool{
	Name:         "claude",
	ConfigDirEnv: "CLAUDE_CONFIG_DIR",
	HomeDir:      ".claude",
}

var nameRe = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

var reservedNames = map[string]bool{
	"init": true, "create": true, "list": true, "ls": true, "rename": true,
	"rm": true, "delete": true, "run": true, "shell": true, "env": true,
	"token": true, "help": true, "version": true, "completion": true,
	"migrate-v2": true,
}

func ValidateProfileName(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}
	if !nameRe.MatchString(name) {
		return fmt.Errorf("invalid profile name %q (allowed: letters, digits, dot, dash, underscore)", name)
	}
	if reservedNames[name] {
		return fmt.Errorf("profile name %q is reserved (it's a subcommand)", name)
	}
	return nil
}

func CemDir() string {
	return filepath.Join(homeDir(), ".config", "cem")
}

func ProfilesDir() string {
	return filepath.Join(CemDir(), "profiles")
}

func ToolProfilesDir() string {
	return filepath.Join(ProfilesDir(), ClaudeTool.Name)
}

func ToolProfileDir(name string) string {
	return filepath.Join(ToolProfilesDir(), name)
}

func ProfilePath(name string) string {
	return filepath.Clean(ToolProfileDir(name))
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
