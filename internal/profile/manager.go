package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type State struct {
	Active string `json:"active"`
}

func ReadState() (*State, error) {
	data, err := os.ReadFile(StatePath())
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("invalid state.json: %w", err)
	}
	return &s, nil
}

func WriteState(s *State) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(StatePath(), data, 0644)
}

func IsInitialized() bool {
	_, err := os.Stat(CemDir())
	return err == nil
}

func Init() error {
	if IsInitialized() {
		return fmt.Errorf("cem is already initialized (found %s)", CemDir())
	}

	profileDir := ProfileDir("default")
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	claudeDir := HomeClaudeDir()
	claudeJSON := HomeClaudeJSON()

	if err := importOrCreate(claudeDir, filepath.Join(profileDir, ".claude"), true); err != nil {
		return err
	}

	if err := importOrCreate(claudeJSON, filepath.Join(profileDir, ".claude.json"), false); err != nil {
		return err
	}

	if err := createSymlinks("default"); err != nil {
		return err
	}

	return WriteState(&State{Active: "default"})
}

func importOrCreate(src, dst string, isDir bool) error {
	info, err := os.Lstat(src)
	if err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(src); err != nil {
				return fmt.Errorf("failed to remove existing symlink %s: %w", src, err)
			}
			return createEmpty(dst, isDir)
		}
		if err := os.Rename(src, dst); err != nil {
			return fmt.Errorf("failed to move %s to %s: %w", src, dst, err)
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat %s: %w", src, err)
	}

	return createEmpty(dst, isDir)
}

func createEmpty(path string, isDir bool) error {
	if isDir {
		return os.MkdirAll(path, 0755)
	}
	return os.WriteFile(path, []byte("{}"), 0644)
}

func Create(name string) error {
	dir := ProfileDir(name)
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("profile %q already exists", name)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, ".claude"), 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, ".claude.json"), []byte("{}"), 0644)
}

func Switch(name string) error {
	if err := CheckClaudeSafe(); err != nil {
		return err
	}

	dir := ProfileDir(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", name)
	}

	for _, path := range []string{HomeClaudeDir(), HomeClaudeJSON()} {
		if err := removeIfSymlink(path); err != nil {
			return err
		}
	}

	if err := createSymlinks(name); err != nil {
		return err
	}

	return WriteState(&State{Active: name})
}

func removeIfSymlink(path string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return os.Remove(path)
	}
	return fmt.Errorf("%s exists but is not a symlink — refusing to remove (run 'cem init' first)", path)
}

func createSymlinks(name string) error {
	dir := ProfileDir(name)

	if err := os.Symlink(filepath.Join(dir, ".claude"), HomeClaudeDir()); err != nil {
		return fmt.Errorf("failed to create symlink for .claude: %w", err)
	}
	if err := os.Symlink(filepath.Join(dir, ".claude.json"), HomeClaudeJSON()); err != nil {
		return fmt.Errorf("failed to create symlink for .claude.json: %w", err)
	}
	return nil
}

func List() ([]string, error) {
	entries, err := os.ReadDir(ProfilesDir())
	if err != nil {
		return nil, fmt.Errorf("failed to read profiles directory: %w", err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func Current() (string, error) {
	state, err := ReadState()
	if err != nil {
		return "", fmt.Errorf("failed to read state: %w", err)
	}
	return state.Active, nil
}

func Rename(oldName, newName string) error {
	oldDir := ProfileDir(oldName)
	newDir := ProfileDir(newName)

	if _, err := os.Stat(oldDir); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", oldName)
	}
	if _, err := os.Stat(newDir); err == nil {
		return fmt.Errorf("profile %q already exists", newName)
	}

	if err := os.Rename(oldDir, newDir); err != nil {
		return fmt.Errorf("failed to rename profile: %w", err)
	}

	state, err := ReadState()
	if err != nil {
		return err
	}
	if state.Active == oldName {
		for _, path := range []string{HomeClaudeDir(), HomeClaudeJSON()} {
			if err := removeIfSymlink(path); err != nil {
				return err
			}
		}
		if err := createSymlinks(newName); err != nil {
			return err
		}
		state.Active = newName
		return WriteState(state)
	}

	return nil
}

func EnsureInitialized() error {
	if !IsInitialized() {
		fmt.Println("cem not initialized. Running first-time setup...")
		return Init()
	}
	return nil
}
