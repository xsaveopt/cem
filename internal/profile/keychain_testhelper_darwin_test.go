//go:build darwin

package profile

import "testing"

type fakeKeychain struct {
	store map[string]string
	calls [][]string
}

func newFakeKeychain() *fakeKeychain {
	return &fakeKeychain{store: map[string]string{}}
}

func (f *fakeKeychain) key(service, account string) string {
	return service + "\x00" + account
}

func (f *fakeKeychain) exec(args ...string) ([]byte, error) {
	f.calls = append(f.calls, args)
	if len(args) == 0 {
		return nil, errFakeNoArgs
	}
	switch args[0] {
	case "find-generic-password":
		s, a := parseSA(args)
		if v, ok := f.store[f.key(s, a)]; ok {
			return []byte(v + "\n"), nil
		}
		return nil, errFakeNotFound
	case "add-generic-password":
		s, a := parseSA(args)
		f.store[f.key(s, a)] = parseW(args)
		return nil, nil
	case "delete-generic-password":
		s, a := parseSA(args)
		if _, ok := f.store[f.key(s, a)]; !ok {
			return nil, errFakeNotFound
		}
		delete(f.store, f.key(s, a))
		return nil, nil
	}
	return nil, errFakeUnknownCmd
}

func parseSA(args []string) (service, account string) {
	for i := 0; i < len(args)-1; i++ {
		switch args[i] {
		case "-s":
			service = args[i+1]
		case "-a":
			account = args[i+1]
		}
	}
	return
}

func parseW(args []string) string {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-w" {
			return args[i+1]
		}
	}
	return ""
}

type fakeErr string

func (e fakeErr) Error() string { return string(e) }

const (
	errFakeNoArgs     fakeErr = "fake keychain: no args"
	errFakeNotFound   fakeErr = "fake keychain: not found"
	errFakeUnknownCmd fakeErr = "fake keychain: unknown subcommand"
)

func installFakeKeychain(t *testing.T) *fakeKeychain {
	t.Helper()
	fk := newFakeKeychain()
	prev := keychainExec
	keychainExec = fk.exec
	t.Cleanup(func() { keychainExec = prev })
	return fk
}
