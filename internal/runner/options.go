package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/pdtm/pkg/path"
	fileutil "github.com/projectdiscovery/utils/file"
	folderutil "github.com/projectdiscovery/utils/folder"
)

var defaultConfigLocation = filepath.Join(folderutil.HomeDirOrDefault("."), ".config/pdtm/config.yaml")
var cacheFile = filepath.Join(folderutil.HomeDirOrDefault("."), ".config/pdtm/cache.json")

// Options contains the configuration options for tuning the enumeration process.
type Options struct {
	ConfigFile string
	Path       string

	Install goflags.StringSlice
	Update  goflags.StringSlice
	Remove  goflags.StringSlice

	InstallAll bool
	UpdateAll  bool
	RemoveAll  bool

	Verbose bool
	Silent  bool
	Version bool
}

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	var err error
	home, err := os.UserHomeDir()
	if err != nil {
		gologger.Fatal().Msgf("Failed to get user home directory: %s", err)
	}
	defaultPath := filepath.Join(home, ".projectdiscovery")
	options := &Options{}
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`projectdiscovery foss tool manager`)

	flagSet.CreateGroup("config", "Config",
		flagSet.StringVar(&options.ConfigFile, "config", defaultConfigLocation, "flag configuration file"),
		flagSet.StringVar(&options.Path, "path", defaultPath, "path"),
	)

	flagSet.CreateGroup("install", "Install",
		flagSet.StringSliceVarP(&options.Install, "install", "i", []string{}, "install given pd-tool (comma separated)", goflags.NormalizedStringSliceOptions),
		flagSet.BoolVarP(&options.InstallAll, "install-all", "ia", false, "install all pd-tools"),
	)

	flagSet.CreateGroup("update", "Update",
		flagSet.StringSliceVarP(&options.Update, "update", "u", []string{}, "update given pd-tool (comma separated)", goflags.NormalizedStringSliceOptions),
		flagSet.BoolVarP(&options.UpdateAll, "update-all", "ua", false, "update all pd-tools"),
	)

	flagSet.CreateGroup("remove", "Remove",
		flagSet.StringSliceVarP(&options.Remove, "remove", "r", []string{}, "remove given pd-tool (comma separated)", goflags.NormalizedStringSliceOptions),
		flagSet.BoolVarP(&options.RemoveAll, "remove-all", "ra", false, "remove all pd-tools"),
	)

	flagSet.CreateGroup("debug", "Debug",
		flagSet.BoolVar(&options.Version, "version", false, "show version of the project"),
		flagSet.BoolVar(&options.Verbose, "v", false, "show verbose output"),
	)

	if err := flagSet.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	options.configureOutput()

	showBanner()

	if options.Version {
		gologger.Info().Msgf("Current Version: %s\n", Version)
		os.Exit(0)
	}

	if options.ConfigFile != defaultConfigLocation {
		_ = options.loadConfigFrom(options.ConfigFile)
	}

	// Validate the options passed by the user and if any
	// invalid options have been used, exit.
	err = options.validateOptions()
	if err != nil {
		gologger.Fatal().Msgf("pdtm error: %s\n", err)
	}

	if options.Path == defaultPath {
		pathVars := strings.Split(os.Getenv("PATH"), ":")
		for _, pathVar := range pathVars {
			if strings.EqualFold(pathVar, defaultPath) {
				return options
			}
		}
		if err := path.SetENV(defaultPath); err != nil {
			gologger.Fatal().Msgf("Failed to set path: %s. Please add ~/.projectdiscovery to $PATH and run again", err)
		}

	}

	return options
}

// configureOutput configures the output on the screen
func (options *Options) configureOutput() {
	// If the user desires verbose output, show verbose output
	if options.Verbose {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)
	}
	if options.Silent {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	}
}

func (Options *Options) loadConfigFrom(location string) error {
	return fileutil.Unmarshal(fileutil.YAML, []byte(location), Options)
}

// validateOptions validates the configuration options passed
func (options *Options) validateOptions() error {
	return nil
}
