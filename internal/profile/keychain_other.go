//go:build !darwin

package profile

func MigrateKeychain(_ string) error          { return nil }
func RenameKeychain(_ string, _ string) error { return nil }
func DeleteKeychain(_ string)                 {}
