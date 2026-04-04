//go:build darwin

package profile

import (
	"fmt"
	"os/exec"
	"os/user"
	"strings"
)

// skipKeychainOps disables keychain operations; set true in tests.
var skipKeychainOps bool

const claudeKeychainService = "Claude Code-credentials"
const cemKeychainService = "cem"

func claudeKeychainAccount() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("could not determine current user: %w", err)
	}
	return u.Username, nil
}

// keychainRead returns the password for the given service+account, or "" if not found.
func keychainRead(service, account string) (string, error) {
	out, err := exec.Command("security", "find-generic-password", "-s", service, "-a", account, "-w").Output()
	if err != nil {
		return "", nil // not found or inaccessible
	}
	return strings.TrimSpace(string(out)), nil
}

// keychainWrite creates or updates a keychain entry.
func keychainWrite(service, account, password string) error {
	err := exec.Command("security", "add-generic-password", "-U", "-s", service, "-a", account, "-w", password).Run()
	if err != nil {
		return fmt.Errorf("keychain write failed: %w", err)
	}
	return nil
}

// keychainDelete removes a keychain entry; ignores errors (including not-found).
func keychainDelete(service, account string) {
	_ = exec.Command("security", "delete-generic-password", "-s", service, "-a", account).Run()
}

// SaveClaudeKeychain backs up the active Claude Code keychain credential under a
// profile-specific slot so it can be restored when switching back to this profile.
func SaveClaudeKeychain(profileName string) error {
	if skipKeychainOps {
		return nil
	}
	account, err := claudeKeychainAccount()
	if err != nil {
		return err
	}
	password, err := keychainRead(claudeKeychainService, account)
	if err != nil {
		return err
	}
	backupAccount := "claude:" + profileName
	if password == "" {
		// No active credential — clear any stale backup for this profile.
		keychainDelete(cemKeychainService, backupAccount)
		return nil
	}
	return keychainWrite(cemKeychainService, backupAccount, password)
}

// RestoreClaudeKeychain loads the stored credential for the target profile and
// makes it the active Claude Code credential. If no backup exists, the active
// credential is deleted (forcing re-authentication).
func RestoreClaudeKeychain(profileName string) error {
	if skipKeychainOps {
		return nil
	}
	account, err := claudeKeychainAccount()
	if err != nil {
		return err
	}
	backupAccount := "claude:" + profileName
	password, err := keychainRead(cemKeychainService, backupAccount)
	if err != nil {
		return err
	}
	if password == "" {
		// No credential stored for this profile — clear the active entry.
		keychainDelete(claudeKeychainService, account)
		return nil
	}
	return keychainWrite(claudeKeychainService, account, password)
}
