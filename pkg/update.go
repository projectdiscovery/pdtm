package pkg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ospath "github.com/projectdiscovery/pdtm/pkg/path"

	"github.com/projectdiscovery/gologger"
)

// Update updates a given tool
func Update(path string, tool Tool) error {
	if executablePath, exists := ospath.GetExecutablePath(path, tool.Name); exists {
		if isUpToDate(tool, path) {
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
		return fmt.Errorf(ErrToolNotFound, tool.Name, executablePath)
	}
}

func isUpToDate(tool Tool, path string) (latest bool) {
	cmd := exec.Command(filepath.Join(path, tool.Name), "--version")

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
