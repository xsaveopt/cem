package profile

import (
	"os"
	"path/filepath"
)

func CemDir() string {
	return filepath.Join(homeDir(), ".config", "cem")
}

func ProfilesDir() string {
	return filepath.Join(CemDir(), "profiles")
}

func ProfileDir(name string) string {
	return filepath.Join(ProfilesDir(), name)
}

func StatePath() string {
	return filepath.Join(CemDir(), "state.json")
}

func HomeClaudeDir() string {
	return filepath.Join(homeDir(), ".claude")
}

func HomeClaudeJSON() string {
	return filepath.Join(homeDir(), ".claude.json")
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
