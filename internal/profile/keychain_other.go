//go:build !darwin

package profile

// skipKeychainOps disables keychain operations; set true in tests.
var skipKeychainOps bool

func SaveClaudeKeychain(_ string) error    { return nil }
func RestoreClaudeKeychain(_ string) error { return nil }
