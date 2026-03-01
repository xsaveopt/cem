package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupTestHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	overrideHome = dir
	skipSafetyCheck = true
	t.Cleanup(func() {
		overrideHome = ""
		skipSafetyCheck = false
	})
	return dir
}

func TestInitFreshNoExistingFiles(t *testing.T) {
	home := setupTestHome(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatalf("ReadState() error: %v", err)
	}
	if state.Active != "default" {
		t.Errorf("state.Active = %q, want %q", state.Active, "default")
	}

	profileDir := ProfileDir("default")
	assertDirExists(t, filepath.Join(profileDir, ".claude"))
	assertFileExists(t, filepath.Join(profileDir, ".claude.json"))

	assertSymlinkTarget(t, filepath.Join(home, ".claude"), filepath.Join(profileDir, ".claude"))
	assertSymlinkTarget(t, filepath.Join(home, ".claude.json"), filepath.Join(profileDir, ".claude.json"))
}

func TestInitImportsExistingFiles(t *testing.T) {
	home := setupTestHome(t)

	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte(`{"key":"val"}`), 0644); err != nil {
		t.Fatal(err)
	}
	claudeJSON := filepath.Join(home, ".claude.json")
	if err := os.WriteFile(claudeJSON, []byte(`{"existing":true}`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	profileDir := ProfileDir("default")
	data, err := os.ReadFile(filepath.Join(profileDir, ".claude", "settings.json"))
	if err != nil {
		t.Fatalf("imported settings.json missing: %v", err)
	}
	if string(data) != `{"key":"val"}` {
		t.Errorf("settings.json content = %q, want %q", string(data), `{"key":"val"}`)
	}

	data, err = os.ReadFile(filepath.Join(profileDir, ".claude.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"existing":true}` {
		t.Errorf(".claude.json content = %q, want %q", string(data), `{"existing":true}`)
	}
}

func TestInitAlreadyInitialized(t *testing.T) {
	setupTestHome(t)

	if err := Init(); err != nil {
		t.Fatalf("first Init() error: %v", err)
	}

	err := Init()
	if err == nil {
		t.Fatal("second Init() should return error")
	}
}

func TestCreateAndList(t *testing.T) {
	setupTestHome(t)
	if err := Init(); err != nil {
		t.Fatal(err)
	}

	if err := Create("work"); err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	names, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("List() returned %d profiles, want 2", len(names))
	}
	if names[0] != "default" || names[1] != "work" {
		t.Errorf("List() = %v, want [default work]", names)
	}
}

func TestCreateDuplicate(t *testing.T) {
	setupTestHome(t)
	if err := Init(); err != nil {
		t.Fatal(err)
	}

	err := Create("default")
	if err == nil {
		t.Fatal("Create(default) should return error for duplicate")
	}
}

func TestSwitchProfile(t *testing.T) {
	home := setupTestHome(t)
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	if err := Create("work"); err != nil {
		t.Fatal(err)
	}

	if err := Switch("work"); err != nil {
		t.Fatalf("Switch() error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.Active != "work" {
		t.Errorf("state.Active = %q, want %q", state.Active, "work")
	}

	workDir := ProfileDir("work")
	assertSymlinkTarget(t, filepath.Join(home, ".claude"), filepath.Join(workDir, ".claude"))
	assertSymlinkTarget(t, filepath.Join(home, ".claude.json"), filepath.Join(workDir, ".claude.json"))
}

func TestSwitchNonexistent(t *testing.T) {
	setupTestHome(t)
	if err := Init(); err != nil {
		t.Fatal(err)
	}

	err := Switch("nope")
	if err == nil {
		t.Fatal("Switch(nope) should return error")
	}
}

func TestRenameActiveProfile(t *testing.T) {
	home := setupTestHome(t)
	if err := Init(); err != nil {
		t.Fatal(err)
	}

	if err := Rename("default", "personal"); err != nil {
		t.Fatalf("Rename() error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.Active != "personal" {
		t.Errorf("state.Active = %q, want %q", state.Active, "personal")
	}

	personalDir := ProfileDir("personal")
	assertSymlinkTarget(t, filepath.Join(home, ".claude"), filepath.Join(personalDir, ".claude"))

	if _, err := os.Stat(ProfileDir("default")); !os.IsNotExist(err) {
		t.Error("old profile dir 'default' should not exist after rename")
	}
}

func TestRenameInactiveProfile(t *testing.T) {
	setupTestHome(t)
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	if err := Create("work"); err != nil {
		t.Fatal(err)
	}

	if err := Rename("work", "office"); err != nil {
		t.Fatalf("Rename() error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.Active != "default" {
		t.Errorf("state.Active = %q, want %q", state.Active, "default")
	}

	assertDirExists(t, ProfileDir("office"))
	if _, err := os.Stat(ProfileDir("work")); !os.IsNotExist(err) {
		t.Error("old profile dir 'work' should not exist after rename")
	}
}

func TestRenameNonexistent(t *testing.T) {
	setupTestHome(t)
	if err := Init(); err != nil {
		t.Fatal(err)
	}

	err := Rename("nope", "other")
	if err == nil {
		t.Fatal("Rename(nope) should return error")
	}
}

func TestRenameToExisting(t *testing.T) {
	setupTestHome(t)
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	if err := Create("work"); err != nil {
		t.Fatal(err)
	}

	err := Rename("work", "default")
	if err == nil {
		t.Fatal("Rename to existing name should return error")
	}
}

func TestCurrent(t *testing.T) {
	setupTestHome(t)
	if err := Init(); err != nil {
		t.Fatal(err)
	}

	name, err := Current()
	if err != nil {
		t.Fatalf("Current() error: %v", err)
	}
	if name != "default" {
		t.Errorf("Current() = %q, want %q", name, "default")
	}
}

func TestStateRoundTrip(t *testing.T) {
	setupTestHome(t)
	if err := os.MkdirAll(CemDir(), 0755); err != nil {
		t.Fatal(err)
	}

	original := &State{Active: "myprofile"}
	if err := WriteState(original); err != nil {
		t.Fatalf("WriteState() error: %v", err)
	}

	loaded, err := ReadState()
	if err != nil {
		t.Fatalf("ReadState() error: %v", err)
	}
	if loaded.Active != original.Active {
		t.Errorf("loaded.Active = %q, want %q", loaded.Active, original.Active)
	}

	data, _ := os.ReadFile(StatePath())
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("state.json is not valid JSON: %v", err)
	}
}

func TestEnsureInitializedRunsInit(t *testing.T) {
	setupTestHome(t)

	if err := EnsureInitialized(); err != nil {
		t.Fatalf("EnsureInitialized() error: %v", err)
	}

	if !IsInitialized() {
		t.Error("should be initialized after EnsureInitialized()")
	}
}

func TestRemoveIfSymlinkRefusesRealFile(t *testing.T) {
	dir := t.TempDir()
	realFile := filepath.Join(dir, "realfile")
	if err := os.WriteFile(realFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	err := removeIfSymlink(realFile)
	if err == nil {
		t.Fatal("removeIfSymlink should refuse to remove a real file")
	}
}

func TestRemoveIfSymlinkHandlesMissing(t *testing.T) {
	err := removeIfSymlink("/nonexistent/path/that/does/not/exist")
	if err != nil {
		t.Fatalf("removeIfSymlink should return nil for missing path, got: %v", err)
	}
}


func assertDirExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("expected dir %s to exist: %v", path, err)
		return
	}
	if !info.IsDir() {
		t.Errorf("expected %s to be a directory", path)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("expected file %s to exist: %v", path, err)
		return
	}
	if info.IsDir() {
		t.Errorf("expected %s to be a file, got directory", path)
	}
}

func assertSymlinkTarget(t *testing.T, link, wantTarget string) {
	t.Helper()
	got, err := os.Readlink(link)
	if err != nil {
		t.Errorf("expected %s to be a symlink: %v", link, err)
		return
	}
	if got != wantTarget {
		t.Errorf("symlink %s -> %q, want %q", link, got, wantTarget)
	}
}
