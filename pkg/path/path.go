package path

import (
	"path/filepath"
	"runtime"
	"strings"

	fileutil "github.com/projectdiscovery/utils/file"
)

type Config struct {
	shellName string
	rcFile    string
}

func SetENV(path string) error {
	_, err := add(path)
	return err
}

func CheckOS() string {
	os := runtime.GOOS
	arc := runtime.GOARCH
	switch os {
	case "windows":
		return "windows_" + arc
	case "darwin":
		return "macOS_" + arc
	case "linux":
		return "linux_" + arc
	default:
		return "not_found"
	}
}

func GetOsData() string {
	os := runtime.GOOS
	arc := runtime.GOARCH
	goVersion := strings.ReplaceAll(runtime.Version(), "go", "")
	return "[OS: " + strings.ToUpper(os) + "] [ARCH: " + strings.ToUpper(arc) + "] [GO: " + goVersion + "]"
}

func GetExecutablePath(path, toolName string) string {
	executablePath := filepath.Join(path, toolName)
	if fileutil.FileExists(executablePath + ".exe") {
		executablePath = executablePath + ".exe"
	}
	return executablePath
}
