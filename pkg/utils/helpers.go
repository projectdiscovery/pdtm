package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/pdtm/pkg"
)

// returns a callback function and when it is executed returns a version string of that tool
func GetVersionCheckCallback(toolName string, output io.Writer) func() {
	return func() {
		tools, err := FetchToolList()
		if err != nil {
			gologger.Error().Msg(err.Error())
			return
		}
		i, exits := Contains(tools, toolName)
		if !exits {
			gologger.Error().Msgf("%v: not found", toolName)
			return
		}
		msg := InstalledVersion(tools[i], au)
		if _, err = output.Write([]byte(fmt.Sprintf("%s %s", toolName, msg))); err != nil {
			gologger.Error().Msg(err.Error())
		}
	}
}

// returns a callback function when executed  updates that tool
func GetUpdaterCallback(toolName string) func() {
	return func() {
		home, _ := os.UserHomeDir()
		dp := filepath.Join(home, ".pdtm/go/bin")
		tools, err := FetchToolList()
		if err != nil {
			gologger.Error().Msg(err.Error())
		}
		i, exits := -1, false
		if i, exits = Contains(tools, toolName); !exits {
			gologger.Error().Msgf("%s doesn't exits", toolName)
			return
		}
		err = pkg.Update(dp, tools[i])
		if err == pkg.ErrIsUpToDate {
			gologger.Info().Msgf("%s: %s", toolName, err)
		} else {
			gologger.Error().Msgf("error while updating %s: %s", toolName, err)
		}
	}
}
