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
	setSkipKeychain(true)
	t.Cleanup(func() {
		overrideHome = ""
		skipSafetyCheck = false
		setSkipKeychain(false)
	})
	return dir
}


func TestInitClaudeFresh(t *testing.T) {
	home := setupTestHome(t)

	if err := Init("claude"); err != nil {
		t.Fatalf("Init(claude) error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatalf("ReadState() error: %v", err)
	}
	if state.Active["claude"] != "default" {
		t.Errorf("state.Active[claude] = %q, want %q", state.Active["claude"], "default")
	}

	profileDir := ToolProfileDir("claude", "default")
	assertDirExists(t, filepath.Join(profileDir, ".claude"))
	assertFileExists(t, filepath.Join(profileDir, ".claude.json"))

	assertSymlinkTarget(t, filepath.Join(home, ".claude"), filepath.Join(profileDir, ".claude"))
	assertSymlinkTarget(t, filepath.Join(home, ".claude.json"), filepath.Join(profileDir, ".claude.json"))
}

func TestInitClaudeImportsExisting(t *testing.T) {
	home := setupTestHome(t)

	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte(`{"key":"val"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".claude.json"), []byte(`{"existing":true}`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Init("claude"); err != nil {
		t.Fatalf("Init(claude) error: %v", err)
	}

	profileDir := ToolProfileDir("claude", "default")

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

	if err := Init("claude"); err != nil {
		t.Fatalf("first Init() error: %v", err)
	}

	err := Init("claude")
	if err == nil {
		t.Fatal("second Init(claude) should return error")
	}
}

func TestCreateAndListClaude(t *testing.T) {
	setupTestHome(t)
	if err := Init("claude"); err != nil {
		t.Fatal(err)
	}

	if err := Create("claude", "work"); err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	workDir := ToolProfileDir("claude", "work")
	assertDirExists(t, filepath.Join(workDir, ".claude"))
	assertFileExists(t, filepath.Join(workDir, ".claude.json"))

	names, err := List("claude")
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
	if err := Init("claude"); err != nil {
		t.Fatal(err)
	}

	err := Create("claude", "default")
	if err == nil {
		t.Fatal("Create(claude, default) should return error for duplicate")
	}
}

func TestSwitchClaudeProfile(t *testing.T) {
	home := setupTestHome(t)
	if err := Init("claude"); err != nil {
		t.Fatal(err)
	}
	if err := Create("claude", "work"); err != nil {
		t.Fatal(err)
	}

	if err := Switch("claude", "work"); err != nil {
		t.Fatalf("Switch() error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.Active["claude"] != "work" {
		t.Errorf("state.Active[claude] = %q, want %q", state.Active["claude"], "work")
	}

	workDir := ToolProfileDir("claude", "work")
	assertSymlinkTarget(t, filepath.Join(home, ".claude"), filepath.Join(workDir, ".claude"))
	assertSymlinkTarget(t, filepath.Join(home, ".claude.json"), filepath.Join(workDir, ".claude.json"))
}

func TestSwitchNonexistent(t *testing.T) {
	setupTestHome(t)
	if err := Init("claude"); err != nil {
		t.Fatal(err)
	}

	err := Switch("claude", "nope")
	if err == nil {
		t.Fatal("Switch(nope) should return error")
	}
}

func TestRenameActiveProfile(t *testing.T) {
	home := setupTestHome(t)
	if err := Init("claude"); err != nil {
		t.Fatal(err)
	}

	if err := Rename("claude", "default", "personal"); err != nil {
		t.Fatalf("Rename() error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.Active["claude"] != "personal" {
		t.Errorf("state.Active[claude] = %q, want %q", state.Active["claude"], "personal")
	}

	personalDir := ToolProfileDir("claude", "personal")
	assertSymlinkTarget(t, filepath.Join(home, ".claude"), filepath.Join(personalDir, ".claude"))

	if _, err := os.Stat(ToolProfileDir("claude", "default")); !os.IsNotExist(err) {
		t.Error("old profile dir 'default' should not exist after rename")
	}
}

func TestRenameInactiveProfile(t *testing.T) {
	setupTestHome(t)
	if err := Init("claude"); err != nil {
		t.Fatal(err)
	}
	if err := Create("claude", "work"); err != nil {
		t.Fatal(err)
	}

	if err := Rename("claude", "work", "office"); err != nil {
		t.Fatalf("Rename() error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.Active["claude"] != "default" {
		t.Errorf("state.Active[claude] = %q, want %q", state.Active["claude"], "default")
	}

	assertDirExists(t, ToolProfileDir("claude", "office"))
	if _, err := os.Stat(ToolProfileDir("claude", "work")); !os.IsNotExist(err) {
		t.Error("old profile dir 'work' should not exist after rename")
	}
}

func TestRenameNonexistent(t *testing.T) {
	setupTestHome(t)
	if err := Init("claude"); err != nil {
		t.Fatal(err)
	}

	err := Rename("claude", "nope", "other")
	if err == nil {
		t.Fatal("Rename(nope) should return error")
	}
}

func TestRenameToExisting(t *testing.T) {
	setupTestHome(t)
	if err := Init("claude"); err != nil {
		t.Fatal(err)
	}
	if err := Create("claude", "work"); err != nil {
		t.Fatal(err)
	}

	err := Rename("claude", "work", "default")
	if err == nil {
		t.Fatal("Rename to existing name should return error")
	}
}

func TestCurrentClaude(t *testing.T) {
	setupTestHome(t)
	if err := Init("claude"); err != nil {
		t.Fatal(err)
	}

	name, err := Current("claude")
	if err != nil {
		t.Fatalf("Current() error: %v", err)
	}
	if name != "default" {
		t.Errorf("Current() = %q, want %q", name, "default")
	}
}


func TestInitGeminiFresh(t *testing.T) {
	home := setupTestHome(t)

	if err := Init("gemini"); err != nil {
		t.Fatalf("Init(gemini) error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatalf("ReadState() error: %v", err)
	}
	if state.Active["gemini"] != "default" {
		t.Errorf("state.Active[gemini] = %q, want %q", state.Active["gemini"], "default")
	}

	profileDir := ToolProfileDir("gemini", "default")
	assertDirExists(t, filepath.Join(profileDir, ".gemini"))
	assertSymlinkTarget(t, filepath.Join(home, ".gemini"), filepath.Join(profileDir, ".gemini"))
}

func TestInitGeminiImportsExisting(t *testing.T) {
	home := setupTestHome(t)

	geminiDir := filepath.Join(home, ".gemini")
	if err := os.MkdirAll(geminiDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(geminiDir, "config.json"), []byte(`{"gemini":true}`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Init("gemini"); err != nil {
		t.Fatalf("Init(gemini) error: %v", err)
	}

	profileDir := ToolProfileDir("gemini", "default")
	data, err := os.ReadFile(filepath.Join(profileDir, ".gemini", "config.json"))
	if err != nil {
		t.Fatalf("imported gemini config.json missing: %v", err)
	}
	if string(data) != `{"gemini":true}` {
		t.Errorf("config.json content = %q, want %q", string(data), `{"gemini":true}`)
	}
}

func TestCreateAndSwitchGemini(t *testing.T) {
	home := setupTestHome(t)
	if err := Init("gemini"); err != nil {
		t.Fatal(err)
	}

	if err := Create("gemini", "work"); err != nil {
		t.Fatalf("Create(gemini, work) error: %v", err)
	}

	workDir := ToolProfileDir("gemini", "work")
	assertDirExists(t, filepath.Join(workDir, ".gemini"))

	if err := Switch("gemini", "work"); err != nil {
		t.Fatalf("Switch(gemini, work) error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.Active["gemini"] != "work" {
		t.Errorf("state.Active[gemini] = %q, want %q", state.Active["gemini"], "work")
	}
	assertSymlinkTarget(t, filepath.Join(home, ".gemini"), filepath.Join(workDir, ".gemini"))
}


func TestInitCopilotFresh(t *testing.T) {
	home := setupTestHome(t)

	if err := Init("copilot"); err != nil {
		t.Fatalf("Init(copilot) error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatalf("ReadState() error: %v", err)
	}
	if state.Active["copilot"] != "default" {
		t.Errorf("state.Active[copilot] = %q, want %q", state.Active["copilot"], "default")
	}

	profileDir := ToolProfileDir("copilot", "default")
	assertDirExists(t, filepath.Join(profileDir, ".copilot"))
	assertSymlinkTarget(t, filepath.Join(home, ".copilot"), filepath.Join(profileDir, ".copilot"))
}

func TestInitCopilotImportsExisting(t *testing.T) {
	home := setupTestHome(t)

	copilotDir := filepath.Join(home, ".copilot")
	if err := os.MkdirAll(copilotDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(copilotDir, "hosts.json"), []byte(`{"copilot":true}`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Init("copilot"); err != nil {
		t.Fatalf("Init(copilot) error: %v", err)
	}

	profileDir := ToolProfileDir("copilot", "default")
	data, err := os.ReadFile(filepath.Join(profileDir, ".copilot", "hosts.json"))
	if err != nil {
		t.Fatalf("imported copilot hosts.json missing: %v", err)
	}
	if string(data) != `{"copilot":true}` {
		t.Errorf("hosts.json content = %q, want %q", string(data), `{"copilot":true}`)
	}
}

func TestCreateAndSwitchCopilot(t *testing.T) {
	home := setupTestHome(t)
	if err := Init("copilot"); err != nil {
		t.Fatal(err)
	}

	if err := Create("copilot", "work"); err != nil {
		t.Fatalf("Create(copilot, work) error: %v", err)
	}

	workDir := ToolProfileDir("copilot", "work")
	assertDirExists(t, filepath.Join(workDir, ".copilot"))

	if err := Switch("copilot", "work"); err != nil {
		t.Fatalf("Switch(copilot, work) error: %v", err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.Active["copilot"] != "work" {
		t.Errorf("state.Active[copilot] = %q, want %q", state.Active["copilot"], "work")
	}
	assertSymlinkTarget(t, filepath.Join(home, ".copilot"), filepath.Join(workDir, ".copilot"))
}


func TestToolsAreIndependent(t *testing.T) {
	home := setupTestHome(t)

	if err := Init("claude"); err != nil {
		t.Fatal(err)
	}
	if err := Init("gemini"); err != nil {
		t.Fatal(err)
	}

	if err := Create("claude", "work"); err != nil {
		t.Fatal(err)
	}
	if err := Switch("claude", "work"); err != nil {
		t.Fatal(err)
	}

	state, err := ReadState()
	if err != nil {
		t.Fatal(err)
	}
	if state.Active["claude"] != "work" {
		t.Errorf("claude active = %q, want %q", state.Active["claude"], "work")
	}
	if state.Active["gemini"] != "default" {
		t.Errorf("gemini active = %q, want %q", state.Active["gemini"], "default")
	}

	claudeWorkDir := ToolProfileDir("claude", "work")
	geminiDefaultDir := ToolProfileDir("gemini", "default")
	assertSymlinkTarget(t, filepath.Join(home, ".claude"), filepath.Join(claudeWorkDir, ".claude"))
	assertSymlinkTarget(t, filepath.Join(home, ".gemini"), filepath.Join(geminiDefaultDir, ".gemini"))
}

func TestInvalidTool(t *testing.T) {
	setupTestHome(t)

	err := Init("vscode")
	if err == nil {
		t.Fatal("Init(vscode) should return error for unknown tool")
	}
}


func TestStateRoundTrip(t *testing.T) {
	setupTestHome(t)
	if err := os.MkdirAll(CemDir(), 0755); err != nil {
		t.Fatal(err)
	}

	original := &State{Active: map[string]string{
		"claude":  "personal",
		"gemini":  "work",
		"copilot": "default",
	}}
	if err := WriteState(original); err != nil {
		t.Fatalf("WriteState() error: %v", err)
	}

	loaded, err := ReadState()
	if err != nil {
		t.Fatalf("ReadState() error: %v", err)
	}
	for tool, want := range original.Active {
		if loaded.Active[tool] != want {
			t.Errorf("loaded.Active[%s] = %q, want %q", tool, loaded.Active[tool], want)
		}
	}

	data, _ := os.ReadFile(StatePath())
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("state.json is not valid JSON: %v", err)
	}
}

func TestEnsureToolInitializedRunsInit(t *testing.T) {
	setupTestHome(t)

	if err := EnsureToolInitialized("claude"); err != nil {
		t.Fatalf("EnsureToolInitialized(claude) error: %v", err)
	}

	if !IsToolInitialized("claude") {
		t.Error("claude should be initialized after EnsureToolInitialized")
	}
}

func TestEnsureToolInitializedPerTool(t *testing.T) {
	setupTestHome(t)

	if err := EnsureToolInitialized("gemini"); err != nil {
		t.Fatalf("EnsureToolInitialized(gemini) error: %v", err)
	}

	if !IsToolInitialized("gemini") {
		t.Error("gemini should be initialized")
	}
	if IsToolInitialized("claude") {
		t.Error("claude should NOT be initialized")
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
