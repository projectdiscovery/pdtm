package pkg

import (
	"os"
	"os/exec"
)

// Remove removes given tool
func Remove(tool Tool) error {
	executablePath, err := exec.LookPath(tool.Name)
	if err != nil {
		return err
	}
	return os.Remove(executablePath)
}
