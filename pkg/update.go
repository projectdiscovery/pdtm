package pkg

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	ospath "github.com/projectdiscovery/pdtm/pkg/path"
	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/projectdiscovery/pdtm/pkg/version"

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
			showReleaseNotes(tool.Repo, version)
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

// showReleaseNotes prints the release body for the version that was actually
// installed. Fetching by tag (instead of "latest") avoids showing notes from
// a release the user did not get, e.g. when api.pdtm.sh returns a cached
// older version. See https://github.com/projectdiscovery/pdtm/issues/435.
func showReleaseNotes(repo, installedVersion string) {
	body, err := fetchReleaseBody(repo, installedVersion)
	if err != nil {
		gologger.Warning().Label("updater").Msgf("could not fetch %s %s release notes: %v", repo, installedVersion, err)
		return
	}
	r, err := glamour.NewTermRenderer(glamour.WithAutoStyle())
	if err != nil {
		gologger.Error().Msgf("markdown rendering not supported: %v", err)
	}
	if rendered, err := r.Render(body); err == nil {
		body = rendered
	} else {
		gologger.Error().Msg(err.Error())
	}
	gologger.Print().Msgf("%v\n", body)
}

func fetchReleaseBody(repo, installedVersion string) (string, error) {
	tag := "v" + strings.TrimPrefix(installedVersion, "v")
	rel, _, err := GithubClient().Repositories.GetReleaseByTag(context.Background(), types.Organization, repo, tag)
	if err != nil {
		return "", err
	}
	return rel.GetBody(), nil
}
