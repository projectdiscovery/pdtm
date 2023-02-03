package path

import (
	"runtime"
	"strings"
)

type Config struct {
	shellName string
	rcFile    string
}

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
