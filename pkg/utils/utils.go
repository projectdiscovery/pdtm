package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v4"
	"github.com/projectdiscovery/pdtm/pkg/path"
	"github.com/projectdiscovery/pdtm/pkg/types"
	"github.com/projectdiscovery/pdtm/pkg/version"
	updateutils "github.com/projectdiscovery/utils/update"
)

var host = getEnv("PDTM_SERVER", "https://api.pdtm.sh")

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// configure aurora for logging
var au = aurora.New(aurora.WithColors(true))

func FetchToolList() ([]types.Tool, error) {
	tools := make([]types.Tool, 0)

	// Create the request URL with query parameters
	reqURL := fmt.Sprintf("%s/api/v1/tools/?%s", host, updateutils.GetpdtmParams(""))

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Just log warning as we're already returning from function
			fmt.Printf("Error closing response body: %s\n", err)
		}
	}()

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

func fetchTool(toolName string) (types.Tool, error) {
	var tool types.Tool
	// Create the request URL to get tool
	reqURL := fmt.Sprintf("%s/api/v1/tools/%s?%s", host, toolName, updateutils.GetpdtmParams(""))
	resp, err := http.Get(reqURL)
	if err != nil {
		return tool, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Just log warning as we're already returning from function
			fmt.Printf("Error closing response body: %s\n", err)
		}
	}()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return tool, err
		}
		// edge case for nuclei coz, the nuclei api send a list of tools including nuclei-templates
		if toolName == "nuclei" {
			var data types.NucleiData
			err = json.Unmarshal(body, &data)
			if err != nil {
				return tool, err
			}
			for _, v := range data.Tools {
				if v.Name == toolName {
					tool = v
					break
				}
			}
			return tool, nil
		}

		err = json.Unmarshal(body, &tool)
		if err != nil {
			return tool, err
		}
		return tool, nil
	}
	return tool, nil
}

func Contains(s []types.Tool, toolName string) (int, bool) {
	for i, a := range s {
		if strings.EqualFold(a.Name, toolName) {
			return i, true
		}
	}
	return -1, false
}

func InstalledVersion(tool types.Tool, basePath string, au *aurora.Aurora) string {
	var msg string

	installedVersion, err := version.ExtractInstalledVersion(tool, basePath)
	if err != nil {
		osAvailable := isOsAvailable(tool)
		if !osAvailable {
			msg = fmt.Sprintf("(%s)", au.Gray(10, "not supported").String())
		} else {
			msg = fmt.Sprintf("(%s)", au.BrightYellow("not installed").String())
		}
	}

	if installedVersion != "" {
		if strings.Contains(tool.Version, installedVersion) {
			msg = fmt.Sprintf("(%s) (%s)", au.BrightGreen("latest").String(), au.BrightGreen(tool.Version).String())
		} else {
			msg = fmt.Sprintf("(%s) (%s) ➡ (%s)",
				au.Red("outdated").String(),
				au.Red(installedVersion).String(),
				au.BrightGreen(tool.Version).String())
		}
	}

	return msg
}

func isOsAvailable(tool types.Tool) bool {
	osData := path.CheckOS()
	for asset := range tool.Assets {
		expectedAssetPrefix := tool.Name + "_" + tool.Version + "_" + osData
		if strings.Contains(asset, expectedAssetPrefix) {
			return true
		}
	}
	return false
}
