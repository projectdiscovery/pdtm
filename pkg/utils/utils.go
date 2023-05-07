package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
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

var regexVersionNumber = regexp.MustCompile(`(?m)[v\s](\d+\.\d+\.\d+)`)

func InstalledVersion(tool pkg.Tool, basePath string, au *aurora.Aurora) string {
	var msg string

	toolPath := filepath.Join(basePath, tool.Name)
	cmd := exec.Command(toolPath, "--version")

	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outb
	err := cmd.Run()
	if err != nil {
		osAvailable := isOsAvailable(tool)
		if !osAvailable {
			msg = fmt.Sprintf("(%s)", au.Gray(10, "not supported").String())
		} else {
			msg = fmt.Sprintf("(%s)", au.BrightYellow("not installed").String())
		}
	}

	if installedVersion := regexVersionNumber.FindString(strings.ToLower(outb.String())); installedVersion != "" {
		installedVersionString := strings.TrimPrefix(strings.TrimSpace(installedVersion), "v")
		if strings.Contains(tool.Version, installedVersionString) {
			msg = fmt.Sprintf("(%s) (%s)", au.BrightGreen("latest").String(), au.BrightGreen(tool.Version).String())
		} else {
			msg = fmt.Sprintf("(%s) (%s) ➡ (%s)",
				au.Red("outdated").String(),
				au.Red(installedVersionString).String(),
				au.BrightGreen(tool.Version).String())
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
