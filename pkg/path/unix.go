//go:build !windows

package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/projectdiscovery/gologger"
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

func add(path string) (bool, error) {
	pathVars := strings.Split(os.Getenv("PATH"), ":")
	for _, pathVar := range pathVars {
		if strings.EqualFold(pathVar, path) {
			return false, nil
		}
	}

	home, err := os.UserHomeDir()
	if nil != err {
		return false, err
	}
	shell := filepath.Base(os.Getenv("SHELL"))
	script := fmt.Sprintf("export PATH=%s:$PATH\n\n", path)
	for _, c := range confList {
		if c.shellName == shell {
			b, err := os.ReadFile(filepath.Join(home, c.rcFile))
			if nil != err {
				return false, err
			}

			lines := strings.Split(strings.TrimSpace(string(b)), "\n")
			for _, line := range lines {
				if strings.EqualFold(line, strings.TrimSpace(script)) {
					gologger.Info().Msgf("Run `source ~/%s` to add %s to $PATH ", c.rcFile, path)
					return true, nil
				}
			}
			f, err := os.OpenFile(filepath.Join(home, c.rcFile), os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return false, err
			}
			script := fmt.Sprintf("\n\n# Generated for pdtm. Do not edit.\n%s", script)
			if _, err := f.Write([]byte(script)); err != nil {
				return false, err
			}
			if err := f.Close(); err != nil {
				return false, err
			}
			gologger.Info().Label("WRN").Msgf("Run `source ~/%s` to add $PATH (%s)", c.rcFile, path)
			return true, nil
		}
	}
	return false, fmt.Errorf("shell not supported, add %s to $PATH env", path)
}
