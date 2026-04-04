package profile

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func IsClaudeRunning() (bool, error) {
	if runtime.GOOS == "windows" {
		out, err := exec.Command("tasklist").Output()
		if err != nil {
			return false, err
		}
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(strings.ToLower(line), "claude.exe") {
				return true, nil
			}
		}
		return false, nil
	}

	// macOS / Linux: use pgrep for precise process-name matching.
	// pgrep exits 1 (with no output) when nothing matches — that's not an error.

	if runtime.GOOS == "darwin" {
		// Claude desktop app bundle
		if out, _ := exec.Command("pgrep", "-f", "Claude.app/Contents/MacOS/Claude").Output(); len(strings.TrimSpace(string(out))) > 0 {
			return true, nil
		}
	}

	// claude CLI binary (exact name)
	if out, _ := exec.Command("pgrep", "-x", "claude").Output(); len(strings.TrimSpace(string(out))) > 0 {
		return true, nil
	}

	return false, nil
}

func HasLockFiles(tool string) (bool, error) {
	t, ok := Tools[tool]
	if !ok {
		return false, ValidateTool(tool)
	}

	lockPatterns := []string{
		"*.lock",
		"*.pid",
		"*.sock",
		"*.socket",
	}

	for _, item := range t.Items {
		if !item.isDir {
			continue
		}
		dir := homePath(item.name)
		for _, pattern := range lockPatterns {
			matches, err := filepath.Glob(filepath.Join(dir, pattern))
			if err != nil {
				return false, err
			}
			if len(matches) > 0 {
				return true, nil
			}

			matches, err = filepath.Glob(filepath.Join(dir, "*", pattern))
			if err != nil {
				return false, err
			}
			if len(matches) > 0 {
				return true, nil
			}
		}
	}

	return false, nil
}

var skipSafetyCheck bool

func CheckSafe(tool string) error {
	if skipSafetyCheck {
		return nil
	}

	if tool == "claude" {
		running, err := IsClaudeRunning()
		if err != nil {
			return fmt.Errorf("failed to check running processes: %w", err)
		}
		if running {
			return fmt.Errorf("claude appears to be running; please close all Claude instances before switching profiles")
		}
	}

	hasLocks, err := HasLockFiles(tool)
	if err != nil {
		return fmt.Errorf("failed to check lock files: %w", err)
	}
	if hasLocks {
		return fmt.Errorf("lock files detected in %s directories; please close all instances before switching profiles", tool)
	}

	return nil
}
