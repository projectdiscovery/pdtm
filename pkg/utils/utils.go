package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/logrusorgru/aurora/v4"
	"github.com/projectdiscovery/pdtm/pkg"
	"github.com/projectdiscovery/pdtm/pkg/path"
)

const host = "https://api.pdtm.sh"

// configure aurora for logging
var au = aurora.New(aurora.WithColors(true))

func FetchToolList() ([]pkg.Tool, error) {
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

func fetchTool(toolName string) (pkg.Tool, error) {
	var tool pkg.Tool
	// Create the request URL to get tool
	reqURL := host + "/api/v1/tools/" + toolName
	resp, err := http.Get(reqURL)
	if err != nil {
		return tool, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return tool, err
		}
		err = json.Unmarshal(body, &tool)
		if err != nil {
			return tool, err
		}
		return tool, nil
	}
	return tool, nil
}

func Contains(s []pkg.Tool, toolName string) (int, bool) {
	for i, a := range s {
		if strings.EqualFold(a.Name, toolName) {
			return i, true
		}
	}
	return -1, false
}

func InstalledVersion(tool pkg.Tool, basePath string, au *aurora.Aurora) string {
	var msg string

	toolPath := filepath.Join(basePath, tool.Name)
	cmd := exec.Command(toolPath, "--version")

	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb
	err := cmd.Run()
	if err != nil {
		var errNotFound *exec.Error
		if errors.As(err, &errNotFound) {
			osAvailable := isOsAvailable(tool)
			if osAvailable {
				msg = au.BrightYellow("(not installed)").String()
			} else {
				msg = au.Gray(10, "(not supported)").String()
			}
		} else {
			msg = "not found"
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
	return msg
}

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
