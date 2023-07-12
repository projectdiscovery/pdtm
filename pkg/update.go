package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	ospath "github.com/projectdiscovery/pdtm/pkg/path"
	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/projectdiscovery/pdtm/pkg/version"
	updateutils "github.com/projectdiscovery/utils/update"

	"github.com/projectdiscovery/gologger"
)

// Update updates a given tool
func Update(path string, tool types.Tool, disableChangeLog bool) error {
	if executablePath, exists := ospath.GetExecutablePath(path, tool.Name); exists {
		if isUpToDate(tool, path) {
			return types.ErrIsUpToDate
		}
		gologger.Info().Msgf("updating %s...", tool.Name)

		if len(tool.Assets) == 0 {
			return fmt.Errorf(types.ErrNoAssetFound, tool.Name, executablePath)
		}

		if err := os.Remove(executablePath); err != nil {
			return err
		}

		version, err := install(tool, path)
		if err != nil {
			return err
		}
		if !disableChangeLog {
			showReleaseNotes(tool.Repo)
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

func showReleaseNotes(toolname string) {
	gh, err := updateutils.NewghReleaseDownloader(toolname)
	if err != nil {
		gologger.Fatal().Label("updater").Msgf("failed to download latest release got %v", err)
	}
	gh.SetToolName(toolname)
	output := gh.Latest.GetBody()
	// adjust colors for both dark / light terminal themes
	r, err := glamour.NewTermRenderer(glamour.WithAutoStyle())
	if err != nil {
		gologger.Error().Msgf("markdown rendering not supported: %v", err)
	}
	if rendered, err := r.Render(output); err == nil {
		output = rendered
	} else {
		gologger.Error().Msg(err.Error())
	}
	gologger.Print().Msgf("%v\n", output)
}
