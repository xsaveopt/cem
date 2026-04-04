//go:build darwin

package profile

func setSkipKeychain(v bool) { skipKeychainOps = v }
