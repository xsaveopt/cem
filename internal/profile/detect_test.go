package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHasLockFilesNone(t *testing.T) {
	home := setupTestHome(t)

		claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	has, err := HasLockFiles()
	if err != nil {
		t.Fatalf("HasLockFiles() error: %v", err)
	}
	if has {
		t.Error("HasLockFiles() = true, want false (no lock files)")
	}
}

func TestHasLockFilesDetectsLock(t *testing.T) {
	home := setupTestHome(t)

	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(claudeDir, "session.lock"), []byte("1234"), 0644); err != nil {
		t.Fatal(err)
	}

	has, err := HasLockFiles()
	if err != nil {
		t.Fatalf("HasLockFiles() error: %v", err)
	}
	if !has {
		t.Error("HasLockFiles() = false, want true")
	}
}

func TestHasLockFilesDetectsPid(t *testing.T) {
	home := setupTestHome(t)

	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(claudeDir, "claude.pid"), []byte("9999"), 0644); err != nil {
		t.Fatal(err)
	}

	has, err := HasLockFiles()
	if err != nil {
		t.Fatalf("HasLockFiles() error: %v", err)
	}
	if !has {
		t.Error("HasLockFiles() = false, want true")
	}
}

func TestHasLockFilesDetectsNestedSocket(t *testing.T) {
	home := setupTestHome(t)

	claudeDir := filepath.Join(home, ".claude")
	subDir := filepath.Join(claudeDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(subDir, "agent.sock"), []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	has, err := HasLockFiles()
	if err != nil {
		t.Fatalf("HasLockFiles() error: %v", err)
	}
	if !has {
		t.Error("HasLockFiles() = false, want true (nested .sock file)")
	}
}

func TestHasLockFilesMissingDir(t *testing.T) {
	setupTestHome(t)
	has, err := HasLockFiles()
	if err != nil {
		t.Fatalf("HasLockFiles() error: %v", err)
	}
	if has {
		t.Error("HasLockFiles() = true, want false (dir doesn't exist)")
	}
}
