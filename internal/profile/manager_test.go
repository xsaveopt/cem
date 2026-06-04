package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	overrideHome = dir
	t.Cleanup(func() { overrideHome = "" })
	return dir
}

func TestValidateProfileName(t *testing.T) {
	cases := []struct {
		name string
		ok   bool
	}{
		{"default", true},
		{"work-1", true},
		{"a.b_c-D", true},
		{"", false},
		{"has space", false},
		{"slash/x", false},
		{"run", false},
		{"ls", false},
		{"create", false},
	}
	for _, c := range cases {
		err := ValidateProfileName(c.name)
		if c.ok && err != nil {
			t.Errorf("ValidateProfileName(%q) err = %v, want nil", c.name, err)
		}
		if !c.ok && err == nil {
			t.Errorf("ValidateProfileName(%q) = nil, want error", c.name)
		}
	}
}

func TestCreateSeedsClaudeJSON(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)

	if err := Create("work"); err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	dir := ToolProfileDir("work")
	assertDirExists(t, dir)
	data, err := os.ReadFile(filepath.Join(dir, ".claude.json"))
	if err != nil {
		t.Fatalf("seed .claude.json missing: %v", err)
	}
	if string(data) != seedClaudeJSON {
		t.Errorf(".claude.json = %q, want %q", string(data), seedClaudeJSON)
	}
}

func TestCreateDuplicate(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)
	if err := Create("work"); err != nil {
		t.Fatal(err)
	}
	if err := Create("work"); err == nil {
		t.Fatal("duplicate Create should error")
	}
}

func TestCreateRejectsInvalidName(t *testing.T) {
	setupTestHome(t)
	if err := Create("run"); err == nil {
		t.Fatal("Create(reserved-name) should error")
	}
	if err := Create("bad/name"); err == nil {
		t.Fatal("Create(invalid-chars) should error")
	}
}

func TestInitFreshNoExistingClaude(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	dir := ToolProfileDir("default")
	assertDirExists(t, dir)
	assertFileExists(t, filepath.Join(dir, ".claude.json"))
}

func TestInitImportsExistingFlattened(t *testing.T) {
	home := setupTestHome(t)
	installFakeKeychain(t)

	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(filepath.Join(claudeDir, "projects"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte(`{"k":"v"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "projects", "p.json"), []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".claude.json"), []byte(`{"existing":true}`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	dir := ToolProfileDir("default")
	got, err := os.ReadFile(filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatalf("settings.json should be at profile root: %v", err)
	}
	if string(got) != `{"k":"v"}` {
		t.Errorf("settings.json = %q", string(got))
	}
	if _, err := os.Stat(filepath.Join(dir, "projects", "p.json")); err != nil {
		t.Errorf("nested file should be preserved: %v", err)
	}
	got, err = os.ReadFile(filepath.Join(dir, ".claude.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"existing":true}` {
		t.Errorf(".claude.json = %q", string(got))
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude")); err == nil {
		t.Error("Init should not create nested .claude/ — that was the v2 layout")
	}

	if _, err := os.Stat(claudeDir); err != nil {
		t.Error("Init should not remove source ~/.claude")
	}
}

func TestInitAlreadyInitialized(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)

	if err := Init(); err != nil {
		t.Fatal(err)
	}
	if err := Init(); err == nil {
		t.Fatal("second Init should error")
	}
}

func TestListEmpty(t *testing.T) {
	setupTestHome(t)
	names, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("List() = %v, want empty", names)
	}
}

func TestListAfterCreate(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)
	for _, n := range []string{"work", "personal", "ci"} {
		if err := Create(n); err != nil {
			t.Fatal(err)
		}
	}
	names, err := List()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"ci", "personal", "work"}
	if len(names) != 3 || names[0] != want[0] || names[1] != want[1] || names[2] != want[2] {
		t.Errorf("List() = %v, want %v", names, want)
	}
}

func TestRename(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)
	if err := Create("work"); err != nil {
		t.Fatal(err)
	}
	if err := Rename("work", "job"); err != nil {
		t.Fatalf("Rename() error: %v", err)
	}
	if Exists("work") {
		t.Error("old profile still exists")
	}
	if !Exists("job") {
		t.Error("new profile missing")
	}
}

func TestRenameNonexistent(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)
	if err := Rename("nope", "other"); err == nil {
		t.Fatal("Rename(nope) should error")
	}
}

func TestRenameToExisting(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)
	if err := Create("a"); err != nil {
		t.Fatal(err)
	}
	if err := Create("b"); err != nil {
		t.Fatal(err)
	}
	if err := Rename("a", "b"); err == nil {
		t.Fatal("Rename to existing should error")
	}
}

func TestDelete(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)
	if err := Create("work"); err != nil {
		t.Fatal(err)
	}
	if err := Delete("work"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
	if Exists("work") {
		t.Error("profile dir still exists after Delete")
	}
}

func TestDeleteNonexistent(t *testing.T) {
	setupTestHome(t)
	if err := Delete("nope"); err == nil {
		t.Fatal("Delete(nope) should error")
	}
}

func TestProfilePathDistinct(t *testing.T) {
	setupTestHome(t)
	if ProfilePath("work") == ProfilePath("personal") {
		t.Error("different names should give different paths")
	}
}

func TestMigrateV2FlattensNestedClaude(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)

	dir := ToolProfileDir("work")
	nested := filepath.Join(dir, ".claude")
	if err := os.MkdirAll(filepath.Join(nested, "projects"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "settings.json"), []byte(`{"v2":1}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "projects", "x.json"), []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"outer":true}`), 0644); err != nil {
		t.Fatal(err)
	}

	rep, err := MigrateV2()
	if err != nil {
		t.Fatalf("MigrateV2: %v", err)
	}
	if len(rep.Flattened) != 1 || rep.Flattened[0] != "work" {
		t.Errorf("Flattened = %v, want [work]", rep.Flattened)
	}

	got, err := os.ReadFile(filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatalf("settings.json should be flattened: %v", err)
	}
	if string(got) != `{"v2":1}` {
		t.Errorf("settings.json = %q", string(got))
	}
	if _, err := os.Stat(filepath.Join(dir, "projects", "x.json")); err != nil {
		t.Errorf("nested file lost: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude")); err == nil {
		t.Error(".claude/ should be removed after flatten")
	}

	got, err = os.ReadFile(filepath.Join(dir, ".claude.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"outer":true}` {
		t.Errorf("outer .claude.json should win conflict, got %q", string(got))
	}
}

func TestMigrateV2Idempotent(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)
	if err := Create("work"); err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateV2(); err != nil {
		t.Fatal(err)
	}
	rep, err := MigrateV2()
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Flattened) != 0 {
		t.Errorf("second run should be a no-op, got %v", rep.Flattened)
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
