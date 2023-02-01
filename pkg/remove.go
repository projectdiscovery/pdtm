package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/projectdiscovery/gologger"
	fileutil "github.com/projectdiscovery/utils/file"
)

// Remove removes given tool
func Remove(path string, tool Tool) error {
	executablePath := filepath.Join(path, tool.Name)
	if fileutil.FileExists(executablePath + ".exe") {
		executablePath = executablePath + ".exe"
	}
	if fileutil.FileExists(executablePath) {
		gologger.Info().Msgf("removing %s...", tool.Name)
		err := os.Remove(executablePath)
		if err != nil {
			return err
		}
		gologger.Info().Msgf("removed %s", tool.Name)
		return nil
	}
	return fmt.Errorf(ErrToolNotFound, tool.Name, executablePath)
}
