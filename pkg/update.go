package pkg

import (
	"bytes"
	"os/exec"
	"strings"
)

// Update updates a given tool
func Update(tool Tool, path string) (string, error) {
	if isUpToDate(tool) {
		return "", ErrIsUpToDate
	}

	if err := Remove(tool); err != nil {
		return "", err
	}
	version, err := Install(tool, path)
	if err != nil {
		return "", err
	}
	return version, nil
}

func isUpToDate(tool Tool) (latest bool) {
	cmd := exec.Command(tool.Name, "--version")

	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb
	err := cmd.Run()
	if err != nil {
		latest = false
		return
	}

	installedVersion := bytes.Split(outb.Bytes(), []byte("Current Version: "))
	if len(installedVersion) == 2 {
		installedVersionString := strings.TrimPrefix(strings.TrimSpace(string(installedVersion[1])), "v")
		if strings.Contains(tool.Version, installedVersionString) {
			latest = true
		} else {
			latest = false
		}
	}
	return
}
