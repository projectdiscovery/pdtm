package pkg

import (
	"os"
	"os/exec"

	stringsutil "github.com/projectdiscovery/utils/strings"

	"github.com/projectdiscovery/gologger"
)

// Remove removes given tool
func Remove(tool Tool) error {
	executablePath, err := exec.LookPath(tool.Name)
	if err != nil {
		return err
	}

	// Check if the path of the executable is in a system directory
	systemDirs := []string{"/usr/bin", "/usr/local/bin", "/bin", "/snap/bin", "/opt/homebrew/bin", "C:/Program Files", "C:/Program Files (x86)"}
	if stringsutil.ContainsAny(executablePath, systemDirs...) {
		gologger.Info().Msgf("skipping removal of %s...", tool.Name)
		return nil
	}

	gologger.Info().Msgf("removing %s...", tool.Name)
	err = os.Remove(executablePath)
	if err != nil {
		return err
	}
	gologger.Info().Msgf("removed %s", tool.Name)
	return nil
}
