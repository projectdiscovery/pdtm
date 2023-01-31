package pkg

import (
	"os"
	"path/filepath"

	"github.com/projectdiscovery/gologger"
)

// Remove removes given tool
func Remove(tool Tool, path string) error {
	executablePath := filepath.Join(path, tool.Name)
	if _, err := os.Stat(executablePath); err == nil {
		gologger.Info().Msgf("removing %s...", tool.Name)
		err = os.Remove(executablePath)
		if err != nil {
			return err
		}
		gologger.Info().Msgf("removed %s", tool.Name)
		return nil
	}
	gologger.Info().Msgf("skipping removal of %s...", tool.Name)
	return nil
}
