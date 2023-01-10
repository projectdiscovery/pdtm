package path

import (
	"runtime"
)

type Config struct {
	shellName string
	rcFile    string
}

func SetENV(path string) error {
	_, err := add(path)
	if err != nil {
		return err
	}
	return nil
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
		return ""
	}
}
