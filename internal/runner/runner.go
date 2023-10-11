package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/pdtm/pkg"
	"github.com/projectdiscovery/pdtm/pkg/path"
	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/projectdiscovery/pdtm/pkg/utils"
	errorutil "github.com/projectdiscovery/utils/errors"
	osutils "github.com/projectdiscovery/utils/os"
	stringsutil "github.com/projectdiscovery/utils/strings"
	"github.com/projectdiscovery/utils/syscallutil"
)

var excludedToolList = []string{"nuclei-templates"}

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
	// add default path to $PATH
	if r.options.SetPath || r.options.Path == defaultPath {
		if err := path.SetENV(r.options.Path); err != nil {
			return errorutil.NewWithErr(err).Msgf(`Failed to set path: %s. Add it to $PATH and run again`, r.options.Path)
		}
	}

	if r.options.UnSetPath {
		if err := path.UnsetENV(r.options.Path); err != nil {
			return errorutil.NewWithErr(err).Msgf(`Failed to unset path: %s. Remove it from $PATH and run again`, r.options.Path)
		}
	}

	toolListApi, err := utils.FetchToolList()
	var toolList []types.Tool

	for _, tool := range toolListApi {
		if !stringsutil.ContainsAny(tool.Name, excludedToolList...) {
			toolList = append(toolList, tool)
		}
	}

	// if toolList is not nil save/update the cache
	// else fetch from cache file
	if toolList != nil {
		go func() {
			if err := UpdateCache(toolList); err != nil {
				gologger.Warning().Msgf("%s\n", err)
			}
		}()
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

	for _, toolName := range r.options.Install {
		if !path.IsSubPath(homeDir, r.options.Path) {
			gologger.Error().Msgf("skipping install outside home folder: %s", toolName)
			continue
		}
		if i, ok := utils.Contains(toolList, toolName); ok {
			tool := toolList[i]
			if tool.InstallType == types.Go && isGoInstalled() {
				if err := pkg.GoInstall(r.options.Path, tool); err != nil {
					gologger.Error().Msgf("%s: %s", tool.Name, err)
				}
				printRequirementInfo(tool)
				continue
			}

			if err := pkg.Install(r.options.Path, tool); err != nil {
				if errors.Is(err, types.ErrIsInstalled) {
					gologger.Info().Msgf("%s: %s", tool.Name, err)
				} else {
					gologger.Error().Msgf("error while installing %s: %s", tool.Name, err)
					gologger.Info().Msgf("trying to install %s using go install", tool.Name)
					if err := pkg.GoInstall(r.options.Path, tool); err != nil {
						gologger.Error().Msgf("%s: %s", tool.Name, err)
					}
				}
			}
			printRequirementInfo(tool)
		} else {
			gologger.Error().Msgf("error while installing %s: %s not found in the list", toolName, toolName)
		}
	}
	for _, tool := range r.options.Update {
		if !path.IsSubPath(homeDir, r.options.Path) {
			gologger.Error().Msgf("skipping update outside home folder: %s", tool)
			continue
		}
		if i, ok := utils.Contains(toolList, tool); ok {
			if err := pkg.Update(r.options.Path, toolList[i], r.options.DisableChangeLog); err != nil {
				if err == types.ErrIsUpToDate {
					gologger.Info().Msgf("%s: %s", tool, err)
				} else {
					gologger.Info().Msgf("%s\n", err)
				}
			}
		}
	}
	for _, tool := range r.options.Remove {
		if !path.IsSubPath(homeDir, r.options.Path) {
			gologger.Error().Msgf("skipping remove outside home folder: %s", tool)
			continue
		}
		if i, ok := utils.Contains(toolList, tool); ok {
			if err := pkg.Remove(r.options.Path, toolList[i]); err != nil {
				var notFoundError *exec.Error
				if errors.As(err, &notFoundError) {
					gologger.Info().Msgf("%s: not found", tool)
				} else {
					gologger.Info().Msgf("%s\n", err)
				}
			}

		}
	}
	if len(r.options.Install) == 0 && len(r.options.Update) == 0 && len(r.options.Remove) == 0 {
		return r.ListToolsAndEnv(toolList)
	}
	return nil
}

func isGoInstalled() bool {
	cmd := exec.Command("go", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func printRequirementInfo(tool types.Tool) {
	specs := getSpecs(tool)

	printTitle := true
	stringBuilder := &strings.Builder{}
	for _, spec := range specs {
		if requirementSatisfied(spec.Name) {
			continue
		}
		if printTitle {
			stringBuilder.WriteString(fmt.Sprintf("%s\n", au.Bold(tool.Name+" requirements:").String()))
			printTitle = false
		}
		instruction := getFormattedInstruction(spec)
		isRequired := getRequirementStatus(spec)
		stringBuilder.WriteString(fmt.Sprintf("%s %s\n", isRequired, instruction))
	}
	if stringBuilder.Len() > 0 {
		gologger.Info().Msgf("%s", stringBuilder.String())
	}
}

func getRequirementStatus(spec types.ToolRequirementSpecification) string {
	if spec.Required {
		return au.Yellow("required").String()
	}
	return au.BrightGreen("optional").String()
}

func getFormattedInstruction(spec types.ToolRequirementSpecification) string {
	return strings.Replace(spec.Instruction, "$CMD", spec.Command, 1)
}

func getSpecs(tool types.Tool) []types.ToolRequirementSpecification {
	var specs []types.ToolRequirementSpecification
	for _, requirement := range tool.Requirements {
		if requirement.OS == runtime.GOOS {
			specs = append(specs, requirement.Specification...)
		}
	}
	return specs
}

// UpdateCache creates/updates cache file
func UpdateCache(toolList []types.Tool) error {
	b, err := json.Marshal(toolList)
	if err != nil {
		return err
	}
	return os.WriteFile(cacheFile, b, os.ModePerm)
}

// FetchFromCache loads tool list from cache file
func FetchFromCache() ([]types.Tool, error) {
	b, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}
	var toolList []types.Tool
	if err := json.Unmarshal(b, &toolList); err != nil {
		return nil, err
	}
	return toolList, nil
}

// ListToolsAndEnv prints the list of tools
func (r *Runner) ListToolsAndEnv(tools []types.Tool) error {
	gologger.Info().Msgf(path.GetOsData() + "\n")
	gologger.Info().Msgf("Path to download project binary: %s\n", r.options.Path)
	var fmtMsg string
	if path.IsSet(r.options.Path) {
		fmtMsg = "Path %s configured in environment variable $PATH\n"
	} else {
		fmtMsg = "Path %s not configured in environment variable $PATH\n"
	}
	gologger.Info().Msgf(fmtMsg, r.options.Path)

	for i, tool := range tools {
		msg := utils.InstalledVersion(tool, r.options.Path, au)
		fmt.Printf("%d. %s %s\n", i+1, tool.Name, msg)
	}
	return nil
}

// Close the runner instance
func (r *Runner) Close() {}

func requirementSatisfied(requirementName string) bool {
	if strings.HasPrefix(requirementName, "lib") {
		libNames := appendLibExtensionForOS(requirementName)
		for _, libName := range libNames {
			_, sysErr := syscallutil.LoadLibrary(libName)
			if sysErr == nil {
				return true
			}
		}
		return false
	}
	_, execErr := exec.LookPath(requirementName)
	return execErr == nil
}

func appendLibExtensionForOS(lib string) []string {
	switch {
	case osutils.IsWindows():
		return []string{fmt.Sprintf("%s.dll", lib), lib}
	case osutils.IsLinux():
		return []string{fmt.Sprintf("%s.so", lib), lib}
	case osutils.IsOSX():
		return []string{fmt.Sprintf("%s.dylib", lib), lib}
	default:
		return []string{lib}
	}
}
