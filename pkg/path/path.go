package path

import (
	osutils "github.com/projectdiscovery/utils/os"
	"path/filepath"
	"runtime"
	"strings"

	fileutil "github.com/projectdiscovery/utils/file"
)

type Config struct {
	shellName string
	rcFile    string
}

func IsSet(path string) bool {
	ok, _ := isSet(path)
	return ok
}

var CommonExtensions = []string{"", ".exe", ".bat"}

func SetENV(path string) error {
	_, err := add(path)
	return err
}

func UnsetENV(path string) error {
	_, err := remove(path)
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

func GetExecutablePath(path, toolName string) (string, bool) {
	basePath := filepath.Join(path, toolName)
	for _, ext := range CommonExtensions {
		executablePath := basePath + ext
		if fileutil.FileExists(executablePath) {
			return executablePath, true
		}
	}

	if osutils.IsWindows() {
		return basePath + ".exe", false
	}

	return basePath, false
}
