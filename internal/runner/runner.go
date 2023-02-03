package runner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/pdtm/pkg"
	"github.com/projectdiscovery/pdtm/pkg/path"
	errorutil "github.com/projectdiscovery/utils/errors"
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
	if r.options.SetPath {
		if err := path.SetENV(r.options.Path); err != nil {
			return errorutil.NewWithErr(err).Msgf(`Failed to set path: %s. Add %s to $PATH and run again`)
		}
	}

	if r.options.UnSetPath {
		if err := path.UnsetENV(r.options.Path); err != nil {
			return errorutil.NewWithErr(err).Msgf(`Failed to set path: %s. Add %s to $PATH and run again`)
		}
	}

	toolList, err := fetchToolList()

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
		if i, ok := contains(toolList, tool); ok {
			if err := pkg.Install(toolList[i], r.options.Path); err != nil {
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
		if i, ok := contains(toolList, tool); ok {
			if err := pkg.Update(toolList[i], r.options.Path); err != nil {
				if err == pkg.ErrIsUpToDate {
					gologger.Info().Msgf("%s: %s", tool, err)
				} else {
					gologger.Error().Msgf("error while updating %s: %s", tool, err)
				}
			}
		}
	}
	for _, tool := range r.options.Remove {
		if i, ok := contains(toolList, tool); ok {
			if err := pkg.Remove(toolList[i]); err != nil {
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
		return r.ListToolsAndEnv(toolList)
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
func (r *Runner) ListToolsAndEnv(tools []pkg.Tool) error {
	gologger.Info().Msgf(path.GetOsData() + "\n")
	gologger.Info().Msgf("Path to download project binary: %s\n\n", r.options.Path)
	paths := path.GetPaths()
	var fmtMsg string
	if path.IsSet(r.options.Path) {
		fmtMsg = "Path %s configured in environment variable $PATH: %s\n"
	} else {
		fmtMsg = "Path %s not configured in environment variable $PATH: %s\n"
	}
	gologger.Info().Msgf(fmtMsg, r.options.Path, strings.Join(paths, ","))
	for i, tool := range tools {
		installedVersion(i+1, tool)
	}
	return nil
}

func installedVersion(i int, tool pkg.Tool) string {
	var msg string

	cmd := exec.Command(tool.Name, "--version")

	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb
	err := cmd.Run()
	if err != nil {
		var notFoundError *exec.Error
		if errors.As(err, &notFoundError) {
			osAvailable := isOsAvailable(tool)
			if osAvailable {
				msg = au.BrightYellow("(not installed)").String()
			} else {
				msg = au.Gray(10, "(not supported)").String()
			}
		} else {
			msg = "version not found"
		}
	}

	installedVersion := strings.Split(strings.ToLower(outb.String()), "current version: ")
	if len(installedVersion) == 2 {
		installedVersionString := strings.TrimPrefix(strings.TrimSpace(string(installedVersion[1])), "v")
		if strings.Contains(tool.Version, installedVersionString) {
			msg = au.Green("(latest) (" + tool.Version + ")").String()
		} else {
			msg = au.Red("(outdated) ("+installedVersionString+")").String() + " ➡ " + au.Green("("+tool.Version+")").String()
		}
	}
	fmt.Printf("%d. %s %s\n", i, tool.Name, msg)
	return msg
}

const host = "https://api.pdtm.sh"

func fetchToolList() ([]pkg.Tool, error) {
	tools := make([]pkg.Tool, 0)

	// Get current OS name, architecture, and Go version
	osName := runtime.GOOS
	osArch := runtime.GOARCH
	goVersion := runtime.Version()

	// Create the request URL with query parameters
	reqURL := host + "/api/v1/tools?os=" + osName + "&arch=" + osArch + "&go_version=" + goVersion

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &tools)
		if err != nil {
			return nil, err
		}
		return tools, nil
	}
	return nil, nil
}

func contains(s []pkg.Tool, toolName string) (int, bool) {
	for i, a := range s {
		if strings.EqualFold(a.Name, toolName) {
			return i, true
		}
	}
	return -1, false
}

// Close the runner instance
func (r *Runner) Close() {}

func isOsAvailable(tool pkg.Tool) bool {
	osData := path.CheckOS()
	for asset := range tool.Assets {
		expectedAssetPrefix := tool.Name + "_" + tool.Version + "_" + osData
		if strings.Contains(asset, expectedAssetPrefix) {
			return true
		}
	}
	return false
}
