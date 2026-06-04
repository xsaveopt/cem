//go:build darwin

package profile

import (
	"strings"
	"testing"
)

func TestClaudeServiceNameDefault(t *testing.T) {
	if got := claudeServiceName(""); got != "Claude Code-credentials" {
		t.Errorf("claudeServiceName(\"\") = %q, want unsuffixed base", got)
	}
}

func TestClaudeServiceNameHashed(t *testing.T) {
	setupTestHome(t)
	a := claudeServiceName("work")
	b := claudeServiceName("personal")

	if !strings.HasPrefix(a, "Claude Code-credentials-") {
		t.Errorf("missing expected prefix: %s", a)
	}
	if len(a) != len("Claude Code-credentials-")+8 {
		t.Errorf("expected 8-char hex suffix: %s", a)
	}
	if a == b {
		t.Error("different profile names should hash to different services")
	}
	if a != claudeServiceName("work") {
		t.Error("hash must be deterministic")
	}
}

func TestMigrateKeychainCopiesEntry(t *testing.T) {
	setupTestHome(t)
	fk := installFakeKeychain(t)

	acct := claudeKeychainAccount()
	fk.store[fk.key("Claude Code-credentials", acct)] = "secret-token"

	if err := MigrateKeychain("default"); err != nil {
		t.Fatalf("MigrateKeychain: %v", err)
	}
	got, ok := fk.store[fk.key(claudeServiceName("default"), acct)]
	if !ok {
		t.Fatal("hashed entry not written")
	}
	if got != "secret-token" {
		t.Errorf("hashed entry = %q, want secret-token", got)
	}
	if fk.store[fk.key("Claude Code-credentials", acct)] != "secret-token" {
		t.Error("source entry should be left in place")
	}
}

func TestMigrateKeychainNoSource(t *testing.T) {
	setupTestHome(t)
	installFakeKeychain(t)
	if err := MigrateKeychain("default"); err != nil {
		t.Errorf("MigrateKeychain with no source should be a no-op, got %v", err)
	}
}

func TestRenameKeychainMovesEntry(t *testing.T) {
	setupTestHome(t)
	fk := installFakeKeychain(t)

	acct := claudeKeychainAccount()
	fk.store[fk.key(claudeServiceName("work"), acct)] = "token-A"

	if err := RenameKeychain("work", "job"); err != nil {
		t.Fatalf("RenameKeychain: %v", err)
	}
	if _, ok := fk.store[fk.key(claudeServiceName("work"), acct)]; ok {
		t.Error("old keychain entry should be deleted")
	}
	if fk.store[fk.key(claudeServiceName("job"), acct)] != "token-A" {
		t.Error("new keychain entry missing or wrong value")
	}
}

func TestMigrateV2KeychainMovesBackup(t *testing.T) {
	setupTestHome(t)
	fk := installFakeKeychain(t)
	if err := Create("work"); err != nil {
		t.Fatal(err)
	}
	fk.store[fk.key("cem", "claude:work")] = "v2-token"

	rep, err := MigrateV2()
	if err != nil {
		t.Fatalf("MigrateV2: %v", err)
	}
	if len(rep.KeychainMoved) != 1 || rep.KeychainMoved[0] != "work" {
		t.Errorf("KeychainMoved = %v, want [work]", rep.KeychainMoved)
	}

	acct := claudeKeychainAccount()
	if got := fk.store[fk.key(claudeServiceName("work"), acct)]; got != "v2-token" {
		t.Errorf("v3 slot = %q, want v2-token", got)
	}
	if _, ok := fk.store[fk.key("cem", "claude:work")]; ok {
		t.Error("v2 backup should be deleted after migration")
	}
}

func TestDeleteKeychainRemovesEntry(t *testing.T) {
	setupTestHome(t)
	fk := installFakeKeychain(t)

	acct := claudeKeychainAccount()
	fk.store[fk.key(claudeServiceName("work"), acct)] = "token"

	DeleteKeychain("work")
	if _, ok := fk.store[fk.key(claudeServiceName("work"), acct)]; ok {
		t.Error("entry should be removed")
	}
}
