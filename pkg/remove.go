package pkg

import (
	"fmt"
	"os"

	ospath "github.com/projectdiscovery/pdtm/pkg/path"
	"github.com/projectdiscovery/pdtm/pkg/types"

	"github.com/projectdiscovery/gologger"
)

// Remove removes given tool
func Remove(path string, tool types.Tool) error {
	executablePath, exists := ospath.GetExecutablePath(path, tool.Name)
	if exists {
		gologger.Info().Msgf("removing %s...", tool.Name)
		err := os.Remove(executablePath)
		if err != nil {
			return err
		}
		gologger.Info().Msgf("removed %s", tool.Name)
		return nil
	}
	return fmt.Errorf(types.ErrToolNotFound, tool.Name, executablePath)
}
