//go:build darwin

package profile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
)

const claudeKeychainBase = "Claude Code-credentials"

var accountRe = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

var keychainExec = func(args ...string) ([]byte, error) {
	return exec.Command("security", args...).Output()
}

func claudeKeychainAccount() string {
	if u := os.Getenv("USER"); accountRe.MatchString(u) {
		return u
	}
	if u, err := user.Current(); err == nil && accountRe.MatchString(u.Username) {
		return u.Username
	}
	return "claude-code-user"
}

func claudeServiceName(name string) string {
	if name == "" {
		return claudeKeychainBase
	}
	sum := sha256.Sum256([]byte(ProfilePath(name)))
	return fmt.Sprintf("%s-%s", claudeKeychainBase, hex.EncodeToString(sum[:])[:8])
}

func keychainRead(service, account string) (string, error) {
	out, err := keychainExec("find-generic-password", "-s", service, "-a", account, "-w")
	if err != nil {
		return "", nil
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func keychainWrite(service, account, password string) error {
	if _, err := keychainExec("add-generic-password", "-U", "-s", service, "-a", account, "-w", password); err != nil {
		return fmt.Errorf("keychain write failed for service %q: %w", service, err)
	}
	return nil
}

func keychainDelete(service, account string) {
	_, _ = keychainExec("delete-generic-password", "-s", service, "-a", account)
}

func MigrateKeychain(name string) error {
	account := claudeKeychainAccount()
	password, err := keychainRead(claudeKeychainBase, account)
	if err != nil {
		return err
	}
	if password == "" {
		return nil
	}
	return keychainWrite(claudeServiceName(name), account, password)
}

func RenameKeychain(oldName, newName string) error {
	account := claudeKeychainAccount()
	password, err := keychainRead(claudeServiceName(oldName), account)
	if err != nil {
		return err
	}
	if password == "" {
		return nil
	}
	if err := keychainWrite(claudeServiceName(newName), account, password); err != nil {
		return err
	}
	keychainDelete(claudeServiceName(oldName), account)
	return nil
}

func DeleteKeychain(name string) {
	keychainDelete(claudeServiceName(name), claudeKeychainAccount())
}
