package pkg

import (
	"fmt"
	"os"
	"strings"

	ospath "github.com/projectdiscovery/pdtm/pkg/path"
	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/projectdiscovery/pdtm/pkg/version"

	"github.com/projectdiscovery/gologger"
)

// Update updates a given tool
func Update(path string, tool types.Tool) error {
	if executablePath, exists := ospath.GetExecutablePath(path, tool.Name); exists {
		if isUpToDate(tool, path) {
			return types.ErrIsUpToDate
		}
		gologger.Info().Msgf("updating %s...", tool.Name)
		if err := os.Remove(executablePath); err != nil {
			return err
		}
		version, err := install(tool, path)
		if err != nil {
			return err
		}
		gologger.Info().Msgf("updated %s to %s (%s)", tool.Name, version, au.BrightGreen("latest").String())
		return nil
	} else {
		return fmt.Errorf(types.ErrToolNotFound, tool.Name, executablePath)
	}
}

func isUpToDate(tool types.Tool, path string) bool {
	v, err := version.ExtractInstalledVersion(tool, path)
	return err == nil && strings.EqualFold(tool.Version, v)
}
