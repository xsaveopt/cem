//go:build !darwin

package profile

func SaveClaudeKeychain(_ string) error    { return nil }
func RestoreClaudeKeychain(_ string) error { return nil }
