package path

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/projectdiscovery/gologger"
)

type Config struct {
	shellName string
	rcFile    string
}

func SetENV(path string) error {
	confList := []*Config{
		{
			shellName: "bash",
			rcFile:    ".bashrc",
		},
		{
			shellName: "zsh",
			rcFile:    ".zshrc",
		},
	}

	home, err := os.UserHomeDir()
	if nil != err {
		return err
	}

	shell := filepath.Base(os.Getenv("SHELL"))
	script := fmt.Sprintf("export PATH=%s:$PATH\n\n", path)

	for _, c := range confList {
		if c.shellName == shell {
			b, err := ioutil.ReadFile(filepath.Join(home, c.rcFile))
			if nil != err {
				return err
			}

			lines := strings.Split(strings.TrimSpace(string(b)), "\n")
			for _, line := range lines {
				if strings.EqualFold(line, strings.TrimSpace(script)) {
					gologger.Info().Msgf("Please run `source ~/%s` or reload terminal to add %s to $PATH ", c.rcFile, path)
					return nil
				}
			}
			f, err := os.OpenFile(filepath.Join(home, c.rcFile), os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			script := fmt.Sprintf("# Generated for pdtm. Do not edit.\n%s", script)
			if _, err := f.Write([]byte(script)); err != nil {
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
			gologger.Info().Msgf("Please run `source ~/%s` or reload terminal to add %s to $PATH ", c.rcFile, path)
		}
	}
	return nil
}
