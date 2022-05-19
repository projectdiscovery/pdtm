package runner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/pdtm/pkg"
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

func (r *Runner) Run() error {

	toolList, err := fetchToolList(r.options.sourceURL)
	if err != nil {
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

	for _, tool := range r.options.Install {
		if t, ok := toolList[tool]; ok {
			gologger.Info().Msgf("trying to install %s", tool)
			if version, err := pkg.Install(t, r.options.Path); err != nil {
				if err == pkg.ErrIsInstalled {
					gologger.Info().Msgf("%s: %s", tool, err)
				} else {
					gologger.Error().Msgf("error while installing %s: %s", tool, err)
				}
			} else {
				gologger.Info().Msgf("Installed %s %s", tool, version)
			}
		}

	}
	for _, tool := range r.options.Update {
		if t, ok := toolList[tool]; ok {
			gologger.Info().Msgf("trying to udpate %s", tool)
			if version, err := pkg.Update(t, r.options.Path); err != nil {
				if err == pkg.ErrIsUpToDate {
					gologger.Info().Msgf("%s: %s", tool, err)
				} else {
					gologger.Error().Msgf("error while updating %s: %s", tool, err)
				}
			} else {
				gologger.Info().Msgf("Updated %s %s", tool, version)
			}
		}
	}
	for _, tool := range r.options.Remove {
		if t, ok := toolList[tool]; ok {
			gologger.Info().Msgf("trying to remove %s", tool)
			if err := pkg.Remove(t); err != nil {
				var notFoundError *exec.Error
				if errors.As(err, &notFoundError) {
					gologger.Info().Msgf("%s: not found", tool)
				} else {
					gologger.Error().Msgf("error while removing %s: %s", tool, err)
				}
			} else {
				gologger.Info().Msgf("Removed %s", tool)
			}

		}
	}
	if len(r.options.Install) == 0 && len(r.options.Update) == 0 && len(r.options.Remove) == 0 {
		return r.ListTools()
	}
	return nil
}

func Contains(s []pkg.Tool, e string) (int, bool) {
	for i, a := range s {
		if strings.EqualFold(a.Name, e) {
			return i, true
		}
	}
	return -1, false
}

//
func (r *Runner) ListTools() error {
	tools, err := fetchToolList(r.options.sourceURL)
	if err != nil {
		gologger.Error().Msgf("error trying to fetch available tool list: %s", err)
		return err
	}
	fmt.Print("Available ProjectDiscovery FOSS Tools\n\n")
	var i int
	for _, tool := range tools {
		i++
		installedVersion(i, tool)
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
			msg = aurora.Red("not installed").String()
		} else {
			msg = "version not found"
		}
	}

	installedVersion := bytes.Split(outb.Bytes(), []byte("Current Version: "))
	if len(installedVersion) == 2 {
		installedVersionString := strings.TrimPrefix(strings.TrimSpace(string(installedVersion[1])), "v")
		if strings.Contains(tool.Version, installedVersionString) {
			msg = aurora.Green("installed - latest").String()
		} else {
			msg = aurora.Yellow("installed - outdated").String()
		}
	}
	fmt.Printf("%d. %s (%s)\n", i, tool.Name, msg)
	return msg
}

func fetchToolList(sourceURL string) (map[string]pkg.Tool, error) {
	tools := make(map[string]pkg.Tool)
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/tools", sourceURL))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &tools); err != nil {
		return nil, err
	}
	return tools, nil
}

// Close the runner instance
func (r *Runner) Close() {}
