package profile

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func IsClaudeRunning() (bool, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	} else {
		cmd = exec.Command("ps", "aux")
	}

	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(string(out), "\n") {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "claude") &&
			!strings.Contains(lower, "cem") &&
			!strings.Contains(lower, "claude.json") {
			if strings.Contains(lower, "/claude") ||
				strings.Contains(lower, "claude ") ||
				strings.Contains(lower, "claude-") {
				return true, nil
			}
		}
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
