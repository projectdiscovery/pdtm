package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/pdtm/pkg"

	"github.com/google/go-github/github"
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

	toolList, err := fetchToolList()
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
		if i, exist := Contains(toolList, tool); exist {
			if version, err := pkg.Install(toolList[i], r.options.Path); err != nil {
				if err != pkg.ErrIsInstalled {
					gologger.Error().Msgf("error while installing %s: %s", tool, err)
				}
			} else {
				gologger.Info().Msgf("Installed %s %s", tool, version)

			}
		}

	}
	for _, tool := range r.options.Update {
		if i, exist := Contains(toolList, tool); exist {
			if version, err := pkg.Update(toolList[i], r.options.Path); err != nil {
				if err != pkg.ErrIsInstalled {
					gologger.Error().Msgf("error while updating %s: %s", tool, err)
				}
			} else {
				gologger.Info().Msgf("Updated %s %s", tool, version)
			}
		}
	}
	for _, tool := range r.options.Remove {
		if i, exist := Contains(toolList, tool); exist {
			if err := pkg.Remove(toolList[i]); err != nil {
				var notFoundError *exec.Error
				if !errors.As(err, &notFoundError) {
					gologger.Error().Msgf("error while removing %s: %s", tool, err)
				}
			} else {
				gologger.Info().Msgf("Removed %s", tool)
			}

		}
	}
	if len(r.options.Install) == 0 && len(r.options.Update) == 0 && len(r.options.Remove) == 0 {
		return ListTools()
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
func ListTools() error {
	var msg string
	tools, err := fetchToolList()
	if err != nil {
		gologger.Error().Msgf("error trying to fetch available tool list: %s", err)
		return err
	}
	fmt.Print("Available ProjectDiscovery FOSS Tools\n\n")

	for i, tool := range tools {
		version, err := fetchLatestVersion(tool.Repo)
		if err != nil {
			gologger.Error().Msgf("error trying to fetch latest version of available tool(%s): %s", tool.Name, err)
			return err
		}
		cmd := exec.Command(tool.Name, "--version")

		var outb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &outb
		err = cmd.Run()
		if err != nil {
			// var notFoundError *exec.Error
			// if !errors.As(err, &notFoundError) {
			// gologger.Error().Msgf("error trying to check installed version of available tool(%s): %s", tool.Name, err)
			// return err
			// }
			msg = "not installed"
		}

		installedVersion := bytes.Split(outb.Bytes(), []byte("Current Version: "))
		if len(installedVersion) == 2 {
			installedVersionString := strings.TrimPrefix(strings.TrimSpace(string(installedVersion[1])), "v")
			if strings.Contains(version, installedVersionString) {
				msg = "latest"
			} else {
				msg = fmt.Sprintf("outdated - %s", bytes.TrimSpace(installedVersion[1]))
			}
		}
		fmt.Printf("%d. %s - %s (%s)\n", i+1, tool.Name, version, msg)
	}
	return nil
}

func fetchToolList() ([]pkg.Tool, error) {
	var tools []pkg.Tool
	data, err := ioutil.ReadFile("tools.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &tools); err != nil {
		return nil, err
	}
	return tools, nil
}

// fetchLatestVersion fetches the latest version of the tool from github api
func fetchLatestVersion(repo string) (string, error) {
	githubClient := pkg.GithubClient()

	releases, _, err := githubClient.Repositories.ListReleases(context.Background(), pkg.Organization, repo, &github.ListOptions{
		PerPage: 1,
	})
	if err != nil {
		return "", err
	}
	if len(releases) == 0 {
		return "", errors.New("could not get latest release")
	}

	return releases[0].GetTagName(), nil
}

// Close the runner instance
func (r *Runner) Close() {}
