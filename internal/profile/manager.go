package profile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

const seedClaudeJSON = `{"hasCompletedOnboarding":true}`

func IsInitialized() bool {
	_, err := os.Stat(ToolProfilesDir())
	return err == nil
}

func Init() error {
	if IsInitialized() {
		return fmt.Errorf("cem already initialized (found %s)", ToolProfilesDir())
	}

	const name = "default"
	if err := ValidateProfileName(name); err != nil {
		return err
	}

	dir := ToolProfileDir(name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	srcDir := filepath.Join(homeDir(), ClaudeTool.HomeDir)
	if err := copyDirContents(srcDir, dir); err != nil {
		return err
	}

	srcJSON := filepath.Join(homeDir(), ".claude.json")
	dstJSON := filepath.Join(dir, ".claude.json")
	if err := copyIfExists(srcJSON, dstJSON); err != nil {
		return err
	}
	if _, err := os.Stat(dstJSON); os.IsNotExist(err) {
		if err := os.WriteFile(dstJSON, []byte(seedClaudeJSON), 0644); err != nil {
			return err
		}
	}

	if err := MigrateKeychain(name); err != nil {
		fmt.Fprintf(os.Stderr, "warning: keychain migration failed: %v\n", err)
	}
	return nil
}

func Create(name string) error {
	if err := ValidateProfileName(name); err != nil {
		return err
	}
	if !IsInitialized() {
		if err := os.MkdirAll(ToolProfilesDir(), 0755); err != nil {
			return err
		}
	}
	dir := ToolProfileDir(name)
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("profile %q already exists", name)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(seedClaudeJSON), 0644)
}

func List() ([]string, error) {
	entries, err := os.ReadDir(ToolProfilesDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
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

func Exists(name string) bool {
	_, err := os.Stat(ToolProfileDir(name))
	return err == nil
}

func Rename(oldName, newName string) error {
	if err := ValidateProfileName(newName); err != nil {
		return err
	}
	oldDir := ToolProfileDir(oldName)
	newDir := ToolProfileDir(newName)

	if _, err := os.Stat(oldDir); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", oldName)
	}
	if _, err := os.Stat(newDir); err == nil {
		return fmt.Errorf("profile %q already exists", newName)
	}

	if err := RenameKeychain(oldName, newName); err != nil {
		fmt.Fprintf(os.Stderr, "warning: keychain rename failed: %v\n", err)
	}

	if err := os.Rename(oldDir, newDir); err != nil {
		return fmt.Errorf("failed to rename profile: %w", err)
	}
	return nil
}

type V2Report struct {
	Flattened     []string
	KeychainMoved []string
	SymlinksFound []string
}

func MigrateV2() (*V2Report, error) {
	rep := &V2Report{}
	names, err := List()
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		dir := ToolProfileDir(name)
		nested := filepath.Join(dir, ClaudeTool.HomeDir)
		if info, err := os.Stat(nested); err == nil && info.IsDir() {
			if err := flattenNestedClaude(dir, nested); err != nil {
				return nil, fmt.Errorf("flatten %s: %w", name, err)
			}
			rep.Flattened = append(rep.Flattened, name)
		}
		if moved, err := migrateV2Keychain(name); err != nil {
			fmt.Fprintf(os.Stderr, "warning: keychain migration failed for %s: %v\n", name, err)
		} else if moved {
			rep.KeychainMoved = append(rep.KeychainMoved, name)
		}
	}
	for _, p := range []string{
		filepath.Join(homeDir(), ClaudeTool.HomeDir),
		filepath.Join(homeDir(), ".claude.json"),
	} {
		if info, err := os.Lstat(p); err == nil && info.Mode()&os.ModeSymlink != 0 {
			rep.SymlinksFound = append(rep.SymlinksFound, p)
		}
	}
	return rep, nil
}

func flattenNestedClaude(dir, nested string) error {
	entries, err := os.ReadDir(nested)
	if err != nil {
		return err
	}
	for _, e := range entries {
		from := filepath.Join(nested, e.Name())
		to := filepath.Join(dir, e.Name())
		if _, err := os.Lstat(to); err == nil {
			if err := os.RemoveAll(from); err != nil {
				return err
			}
			continue
		}
		if err := os.Rename(from, to); err != nil {
			return err
		}
	}
	return os.Remove(nested)
}

func Delete(name string) error {
	dir := ToolProfileDir(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", name)
	}
	DeleteKeychain(name)
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to remove profile directory: %w", err)
	}
	return nil
}

func copyDirContents(src, dst string) error {
	entries, err := os.ReadDir(src)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	for _, e := range entries {
		from := filepath.Join(src, e.Name())
		to := filepath.Join(dst, e.Name())
		if _, err := os.Lstat(to); err == nil {
			continue
		}
		if err := copyIfExists(from, to); err != nil {
			return err
		}
	}
	return nil
}

func copyIfExists(src, dst string) error {
	info, err := os.Lstat(src)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	if info.IsDir() {
		return copyTree(src, dst)
	}
	return copyFile(src, dst, info.Mode())
}

func copyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, target)
		}
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
