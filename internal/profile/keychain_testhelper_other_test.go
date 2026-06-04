//go:build !darwin

package profile

import "testing"

type fakeKeychain struct{ store map[string]string }

func installFakeKeychain(_ *testing.T) *fakeKeychain {
	return &fakeKeychain{store: map[string]string{}}
}
