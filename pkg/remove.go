package pkg

import (
	"os"
	"os/exec"

	"github.com/projectdiscovery/gologger"
)

// Remove removes given tool
func Remove(tool Tool) error {
	executablePath, err := exec.LookPath(tool.Name)
	if err != nil {
		return err
	}
	gologger.Info().Msgf("removing %s...", tool.Name)

	err = os.Remove(executablePath)
	if err != nil {
		return err
	}
	gologger.Info().Msgf("removed %s", tool.Name)
	return nil
}
