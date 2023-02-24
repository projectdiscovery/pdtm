package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/pdtm/pkg"
	"github.com/projectdiscovery/pdtm/pkg/path"
	"github.com/projectdiscovery/pdtm/pkg/utils"
)

// Runner contains the internal logic of the program
type Runner struct {
	options *Options
}

// NewRunner instance
func NewRunner(options *Options) (*Runner, error) {
	return &Runner{
		options: options,
	}, nil
}

// Run the instance
func (r *Runner) Run() error {
	toolList, err := utils.FetchToolList()

	// if toolList is not nil save/update the cache
	// else fetch from cache file
	if toolList != nil {
		go UpdateCache(toolList) //nolint:errcheck
	} else {
		toolList, err = FetchFromCache()
		if err != nil {
			return errors.New("pdtm api is down, please try again later")
		}
		if toolList != nil {
			gologger.Warning().Msg("pdtm api is down, using cached information while we fix the issue \n\n")
		}
	}
	if toolList == nil && err != nil {
		return err
	}

	switch {
	case r.options.InstallAll:
		for _, tool := range toolList {
			r.options.Install = append(r.options.Install, tool.Name)
		}
	case r.options.UpdateAll:
		for _, tool := range toolList {
			r.options.Update = append(r.options.Update, tool.Name)
		}
	case r.options.RemoveAll:
		for _, tool := range toolList {
			r.options.Remove = append(r.options.Remove, tool.Name)
		}
	}
	gologger.Verbose().Msgf("using path %s", r.options.Path)

	for _, tool := range r.options.Install {
		if i, ok := utils.Contains(toolList, tool); ok {
			if err := pkg.Install(r.options.Path, toolList[i]); err != nil {
				if err == pkg.ErrIsInstalled {
					gologger.Info().Msgf("%s: %s", tool, err)
				} else {
					gologger.Error().Msgf("error while installing %s: %s", tool, err)
				}
			}
		} else {
			gologger.Error().Msgf("error while installing %s: %s not found in the list", tool, tool)
		}
	}
	for _, tool := range r.options.Update {
		if i, ok := utils.Contains(toolList, tool); ok {
			if err := pkg.Update(r.options.Path, toolList[i]); err != nil {
				if err == pkg.ErrIsUpToDate {
					gologger.Info().Msgf("%s: %s", tool, err)
				} else {
					gologger.Error().Msgf("error while updating %s: %s", tool, err)
				}
			}
		}
	}
	for _, tool := range r.options.Remove {
		if i, ok := utils.Contains(toolList, tool); ok {
			if err := pkg.Remove(r.options.Path, toolList[i]); err != nil {
				var notFoundError *exec.Error
				if errors.As(err, &notFoundError) {
					gologger.Info().Msgf("%s: not found", tool)
				} else {
					gologger.Error().Msgf("error while removing %s: %s", tool, err)
				}
			}

		}
	}
	if len(r.options.Install) == 0 && len(r.options.Update) == 0 && len(r.options.Remove) == 0 {
		return r.ListTools(toolList)
	}
	return nil
}

// UpdateCache creates/updates cache file
func UpdateCache(toolList []pkg.Tool) error {
	b, err := json.Marshal(toolList)
	if err != nil {
		return err
	}
	return os.WriteFile(cacheFile, b, os.ModePerm)
}

// FetchFromCache loads tool list from cache file
func FetchFromCache() ([]pkg.Tool, error) {
	b, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}
	var toolList []pkg.Tool
	if err := json.Unmarshal(b, &toolList); err != nil {
		return nil, err
	}
	return toolList, nil
}

// ListTools prints the list of tools
func (r *Runner) ListTools(tools []pkg.Tool) error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	gologger.Info().Msgf(path.GetOsData() + "\n")
	gologger.Info().Msgf("Path to download project binary: %s/.pdtm/go/bin\n\n", dirname)
	var i int
	for _, tool := range tools {
		i++
		msg := utils.InstalledVersion(tool, au)
		fmt.Printf("%d. %s %s\n", i, tool.Name, msg)
	}
	return nil
}

// Close the runner instance
func (r *Runner) Close() {}
