package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)


type State struct {
	Active map[string]string `json:"active"`
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
	if s.Active == nil {
		s.Active = make(map[string]string)
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

func IsToolInitialized(tool string) bool {
	_, err := os.Stat(ToolProfilesDir(tool))
	return err == nil
}

func Init(tool string) error {
	t, ok := Tools[tool]
	if !ok {
		return ValidateTool(tool)
	}

	if IsToolInitialized(tool) {
		return fmt.Errorf("%s profiles already initialized (found %s)", tool, ToolProfilesDir(tool))
	}

	profileDir := ToolProfileDir(tool, "default")
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	for _, item := range t.Items {
		src := homePath(item.name)
		dst := filepath.Join(profileDir, item.name)
		if err := importOrCreate(src, dst, item.isDir); err != nil {
			return err
		}
	}

	if err := createSymlinks(tool, "default"); err != nil {
		return err
	}

	state := &State{Active: make(map[string]string)}
	if IsInitialized() {
		existing, err := ReadState()
		if err == nil {
			state = existing
		}
	}
	state.Active[tool] = "default"
	return WriteState(state)
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

func Create(tool, name string) error {
	t, ok := Tools[tool]
	if !ok {
		return ValidateTool(tool)
	}

	dir := ToolProfileDir(tool, name)
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("%s profile %q already exists", tool, name)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	for _, item := range t.Items {
		path := filepath.Join(dir, item.name)
		if err := createEmpty(path, item.isDir); err != nil {
			return err
		}
	}

	return nil
}

func Switch(tool, name string) error {
	t, ok := Tools[tool]
	if !ok {
		return ValidateTool(tool)
	}

	if err := CheckSafe(tool); err != nil {
		return err
	}

	dir := ToolProfileDir(tool, name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("%s profile %q does not exist", tool, name)
	}

	state, err := ReadState()
	if err != nil {
		return err
	}
	currentProfile := state.Active[tool]

	// Back up the current profile's keychain credential before switching.
	if tool == "claude" && currentProfile != "" && currentProfile != name {
		if err := SaveClaudeKeychain(currentProfile); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to save keychain credentials for profile %q: %v\n", currentProfile, err)
		}
	}

	for _, item := range t.Items {
		if err := removeIfSymlink(homePath(item.name)); err != nil {
			return err
		}
	}

	if err := createSymlinks(tool, name); err != nil {
		return err
	}

	// Restore the target profile's keychain credential.
	if tool == "claude" && currentProfile != name {
		if err := RestoreClaudeKeychain(name); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to restore keychain credentials for profile %q: %v\n", name, err)
		}
	}

	state.Active[tool] = name
	return WriteState(state)
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
	return fmt.Errorf("%s exists but is not a symlink — refusing to remove (run 'cem init --tool <tool>' first)", path)
}

func createSymlinks(tool, name string) error {
	t := Tools[tool]
	dir := ToolProfileDir(tool, name)

	for _, item := range t.Items {
		target := filepath.Join(dir, item.name)
		link := homePath(item.name)
		if err := os.Symlink(target, link); err != nil {
			return fmt.Errorf("failed to create symlink for %s: %w", item.name, err)
		}
	}
	return nil
}

func List(tool string) ([]string, error) {
	if err := ValidateTool(tool); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(ToolProfilesDir(tool))
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

func Current(tool string) (string, error) {
	if err := ValidateTool(tool); err != nil {
		return "", err
	}

	state, err := ReadState()
	if err != nil {
		return "", fmt.Errorf("failed to read state: %w", err)
	}
	active, ok := state.Active[tool]
	if !ok {
		return "", fmt.Errorf("no active profile for %s (run 'cem init --tool %s' first)", tool, tool)
	}
	return active, nil
}

func Rename(tool, oldName, newName string) error {
	t, ok := Tools[tool]
	if !ok {
		return ValidateTool(tool)
	}

	oldDir := ToolProfileDir(tool, oldName)
	newDir := ToolProfileDir(tool, newName)

	if _, err := os.Stat(oldDir); os.IsNotExist(err) {
		return fmt.Errorf("%s profile %q does not exist", tool, oldName)
	}
	if _, err := os.Stat(newDir); err == nil {
		return fmt.Errorf("%s profile %q already exists", tool, newName)
	}

	if err := os.Rename(oldDir, newDir); err != nil {
		return fmt.Errorf("failed to rename profile: %w", err)
	}

	state, err := ReadState()
	if err != nil {
		return err
	}
	if state.Active[tool] == oldName {
		for _, item := range t.Items {
			if err := removeIfSymlink(homePath(item.name)); err != nil {
				return err
			}
		}
		if err := createSymlinks(tool, newName); err != nil {
			return err
		}
		state.Active[tool] = newName
		return WriteState(state)
	}

	return nil
}

func EnsureToolInitialized(tool string) error {
	if !IsToolInitialized(tool) {
		fmt.Printf("cem not initialized for %s. Running first-time setup...\n", tool)
		return Init(tool)
	}
	return nil
}
