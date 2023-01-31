package pkg

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/projectdiscovery/gologger"
	fileutil "github.com/projectdiscovery/utils/file"
)

// Update updates a given tool
func Update(tool Tool, path string) error {
	executablePath := filepath.Join(path, tool.Name)
	if fileutil.FileExists(executablePath) {
		if isUpToDate(tool) {
			return ErrIsUpToDate
		}
		gologger.Info().Msgf("updating %s...", tool.Name)
		if err := os.Remove(executablePath); err != nil {
			return err
		}
		version, err := install(tool, path)
		if err != nil {
			return err
		}
		gologger.Info().Msgf("updated %s to %s(latest)", tool.Name, version)
		return nil
	} else {
		gologger.Info().Msgf("%s: not found in path %s", tool.Name, executablePath)
		return Install(tool, path)
	}
}

func isUpToDate(tool Tool) (latest bool) {
	cmd := exec.Command(tool.Name, "--version")

	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb
	err := cmd.Run()
	if err != nil {
		latest = false
		return
	}

	installedVersion := bytes.Split(outb.Bytes(), []byte("Current Version: "))
	if len(installedVersion) == 2 {
		installedVersionString := strings.TrimPrefix(strings.TrimSpace(string(installedVersion[1])), "v")
		if strings.Contains(tool.Version, installedVersionString) {
			latest = true
		} else {
			latest = false
		}
	}
	return
}
