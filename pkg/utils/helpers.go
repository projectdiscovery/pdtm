package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/pdtm/pkg"
)

// GetVersionCheckCallback returns a callback function and when it is executed returns a version string of that tool
func GetVersionCheckCallback(toolName, basePath string) func() string {
	return func() string {
		tool, err := fetchTool(toolName)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("%s %s", toolName, InstalledVersion(tool, basePath, au))
	}
}

// GetUpdaterCallback returns a callback function when executed  updates that tool
func GetUpdaterCallback(toolName string) func() {
	return func() {
		home, _ := os.UserHomeDir()
		dp := filepath.Join(home, ".pdtm/go/bin")
		tool, err := fetchTool(toolName)
		if err != nil {
			gologger.Error().Msg(err.Error())
		}
		err = pkg.Update(dp, tool)
		if err == pkg.ErrIsUpToDate {
			gologger.Info().Msgf("%s: %s", toolName, err)
		} else {
			gologger.Error().Msgf("error while updating %s: %s", toolName, err)
		}
	}
}
