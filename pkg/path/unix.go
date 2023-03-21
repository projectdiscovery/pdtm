//go:build !windows

package path

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/projectdiscovery/gologger"
	errorutil "github.com/projectdiscovery/utils/errors"
	fileutil "github.com/projectdiscovery/utils/file"
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

func (c *Config) GetRCFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if nil != err {
		return "", err
	}
	rcFilePath := filepath.Join(home, c.rcFile)
	if fileutil.FileExists(rcFilePath) {
		return rcFilePath, nil
	}
	// if file doesn't exist create empty file
	if err := os.WriteFile(rcFilePath, []byte("#\n"), 0644); err != nil {
		return "", fmt.Errorf("failed to create rcFile %v got %v", rcFilePath, err)
	}
	return rcFilePath, nil
}

func lookupConfFromShell() (*Config, error) {
	shell := filepath.Base(os.Getenv("SHELL"))
	for _, conf := range confList {
		if conf.shellName == shell {
			if _, err := conf.GetRCFilePath(); err != nil {
				return nil, err
			}
			return conf, nil
		}
	}
	// assume bash as default shell if variable is empty in unix distros
	if shell == "." && len(confList) > 1 {
		conf := confList[0]
		if _, err := conf.GetRCFilePath(); err != nil {
			return nil, err
		}
		return conf, nil
	}
	return nil, errors.New("shell not supported")
}

func isSet(path string) (bool, error) {
	pathVars := paths()
	return sliceutil.Contains(pathVars, path), nil
}

func add(path string) (bool, error) {
	pathVars := paths()
	if sliceutil.Contains(pathVars, path) {
		return false, nil
	}

	conf, err := lookupConfFromShell()
	if err != nil {
		return false, errorutil.NewWithErr(err).Msgf("add %s to $PATH env", path)
	}

	script := fmt.Sprintf("export PATH=$PATH:%s\n\n", path)
	return exportToConfig(conf, path, script)
}

func remove(path string) (bool, error) {
	pathVars := paths()
	if !sliceutil.Contains(pathVars, path) {
		return false, nil
	}

	conf, err := lookupConfFromShell()
	if err != nil {
		return false, errorutil.NewWithErr(err).Msgf("remove %s from $PATH env", path)
	}
	pathVars = sliceutil.PruneEqual(pathVars, path)
	script := fmt.Sprintf("export PATH=%s\n\n", strings.Join(pathVars, ":"))
	return exportToConfig(conf, path, script)
}

func paths() []string {
	return strings.Split(os.Getenv("PATH"), ":")
}

func exportToConfig(config *Config, path, script string) (bool, error) {
	rcFilePath, err := config.GetRCFilePath()
	if err != nil {
		return false, err
	}
	b, err := os.ReadFile(rcFilePath)
	if err != nil {
		return false, err
	}

	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	for _, line := range lines {
		if strings.EqualFold(line, strings.TrimSpace(script)) {
			gologger.Info().Msgf("Run `source ~/%s` to add %s to $PATH ", config.rcFile, path)
			return true, nil
		}
	}
	f, err := os.OpenFile(rcFilePath, os.O_APPEND|os.O_WRONLY, 0644)
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
