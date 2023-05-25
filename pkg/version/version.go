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
	cmd := exec.Command(toolPath, "--version")

	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	if installedVersion := RegexVersionNumber.FindString(strings.ToLower(outb.String())); installedVersion != "" {
		installedVersionString := strings.TrimPrefix(strings.TrimSpace(installedVersion), "v")
		return installedVersionString, nil
	}

	return "", errors.New("unable to extract installed version")
}
