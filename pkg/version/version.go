package version

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/projectdiscovery/pdtm/pkg/types"
)

var (
	RegexVersionNumber = regexp.MustCompile(`(?m)[v\s](\d+\.\d+\.\d+)`)
	versionCommands    = []string{"--version", "version"}

	// ErrNotInstalled is returned when the tool binary is missing.
	ErrNotInstalled = errors.New("not installed")
	// ErrVersionUnknown is returned when the binary exists but its version
	// could not be determined (unsupported flag, empty/unparseable output, ...).
	// Callers should surface a generic label and may log the wrapped detail.
	ErrVersionUnknown = errors.New("unknown")
)

func ExtractInstalledVersion(tool types.Tool, basePath string) (string, error) {
	if err := types.ValidateToolName(tool.Name); err != nil {
		return "", err
	}
	toolPath := filepath.Join(basePath, tool.Name)

	var lastErr error
	for _, versionCmd := range versionCommands {
		ver, err := tryVersionCommand(toolPath, versionCmd)
		if err == nil {
			return ver, nil
		}
		// A missing binary won't be resolved by trying another command, and
		// bailing keeps the surfaced error independent of command ordering.
		if errors.Is(err, ErrNotInstalled) {
			return "", err
		}
		lastErr = err
	}

	return "", lastErr
}

func tryVersionCommand(toolPath, versionCmd string) (string, error) {
	cmd := exec.Command(toolPath, versionCmd)
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb

	if err := cmd.Run(); err != nil {
		// A missing binary surfaces as fs.ErrNotExist on unix and as
		// exec.ErrNotFound on windows (the .exe is not found on PATH), so both
		// must be treated as "not installed".
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, exec.ErrNotFound) {
			return "", ErrNotInstalled
		}
		return "", fmt.Errorf("%w: %q failed: %v", ErrVersionUnknown, versionCmd, err)
	}

	output := outb.String()
	if output == "" {
		return "", fmt.Errorf("%w: %q produced no output", ErrVersionUnknown, versionCmd)
	}

	installedVersion := RegexVersionNumber.FindString(strings.ToLower(output))
	if installedVersion == "" {
		return "", fmt.Errorf("%w: no version found in %q output", ErrVersionUnknown, versionCmd)
	}

	ver := strings.TrimSpace(installedVersion)
	ver = strings.TrimPrefix(ver, "v")

	return ver, nil
}
