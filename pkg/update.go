package pkg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ospath "github.com/projectdiscovery/pdtm/pkg/path"
	stringsutil "github.com/projectdiscovery/utils/strings"

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

func isUpToDate(tool Tool, path string) bool {
	cmd := exec.Command(filepath.Join(path, tool.Name), "--version")

	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb
	err := cmd.Run()
	if err != nil {
		return false
	}
	out := strings.ToLower(outb.String())
	v, err := stringsutil.Between(out, "current version:", "\n")
	v = strings.Trim(v, "\n v")
	return err == nil && strings.EqualFold(tool.Version, v)
}
