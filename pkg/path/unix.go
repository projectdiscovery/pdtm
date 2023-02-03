//go:build !windows

package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/projectdiscovery/gologger"
	sliceutil "github.com/projectdiscovery/utils/slice"
)

var confList = []*Config{
	{
		shellName: "bash",
		rcFile:    ".bashrc",
	},
	{
		shellName: "zsh",
		rcFile:    ".zshrc",
	},
}

func lookupConfFromShell() (*Config, bool) {
	shell := filepath.Base(os.Getenv("SHELL"))
	for _, conf := range confList {
		if conf.shellName == shell {
			return conf, true
		}
	}
	return nil, false
}

func isSet(path string) (bool, error) {
	pathVars := getPathsFromEnv()
	return sliceutil.Contains(pathVars, path), nil
}

func add(path string) (bool, error) {
	pathVars := getPathsFromEnv()
	if sliceutil.Contains(pathVars, path) {
		return false, nil
	}

	conf, ok := lookupConfFromShell()
	if !ok {
		return false, fmt.Errorf("shell not supported, add %s to $PATH env", path)
	}

	script := fmt.Sprintf("export PATH=$PATH:%s\n\n", path)
	return exportToConfig(conf, path, script)
}

func remove(path string) (bool, error) {
	pathVars := getPathsFromEnv()
	if !sliceutil.Contains(pathVars, path) {
		return false, nil
	}

	conf, ok := lookupConfFromShell()
	if !ok {
		return false, fmt.Errorf("shell not supported, add %s to $PATH env", path)
	}
	pathVars = sliceutil.PruneEqual(pathVars, path)
	script := fmt.Sprintf("export PATH=%s\n\n", strings.Join(pathVars, ":"))
	return exportToConfig(conf, path, script)
}

func getPathsFromEnv() []string {
	return strings.Split(os.Getenv("PATH"), ":")
}

func exportToConfig(config *Config, path, script string) (bool, error) {
	home, err := os.UserHomeDir()
	if nil != err {
		return false, err
	}
	b, err := os.ReadFile(filepath.Join(home, config.rcFile))
	if nil != err {
		return false, err
	}

	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	for _, line := range lines {
		if strings.EqualFold(line, strings.TrimSpace(script)) {
			gologger.Info().Msgf("Run `source ~/%s` to add %s to $PATH ", config.rcFile, path)
			return true, nil
		}
	}
	f, err := os.OpenFile(filepath.Join(home, config.rcFile), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return false, err
	}
	script = fmt.Sprintf("\n\n# Generated for pdtm. Do not edit.\n%s", script)
	if _, err := f.Write([]byte(script)); err != nil {
		return false, err
	}
	if err := f.Close(); err != nil {
		return false, err
	}
	gologger.Info().Label("WRN").Msgf("Run `source ~/%s` to add $PATH (%s)", config.rcFile, path)
	return true, nil
}
