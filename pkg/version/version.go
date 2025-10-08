package version

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/projectdiscovery/pdtm/pkg/types"
)

var RegexVersionNumber = regexp.MustCompile(`(?m)[v\s](\d+\.\d+\.\d+)`)

func ExtractInstalledVersion(tool types.Tool, basePath string) (string, error) {
	toolPath := filepath.Join(basePath, tool.Name)

	versionCommands := []string{"--version", "version"}

	for _, versionCmd := range versionCommands {
		if version, err := tryVersionCommand(toolPath, versionCmd); err == nil {
			return version, nil
		}
	}

	return "", errors.New("unable to extract installed version")
}

func tryVersionCommand(toolPath, versionCmd string) (string, error) {
	cmd := exec.Command(toolPath, versionCmd)
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb

	if err := cmd.Run(); err != nil {
		return "", err
	}

	output := outb.String()
	if output == "" {
		return "", errors.New("empty output")
	}

	installedVersion := RegexVersionNumber.FindString(strings.ToLower(output))
	if installedVersion == "" {
		return "", errors.New("no version found in output")
	}

	version := strings.TrimSpace(installedVersion)
	version = strings.TrimPrefix(version, "v")

	return version, nil
}
