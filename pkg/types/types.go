package types

import (
	"errors"
	"fmt"
)

const Organization = "projectdiscovery"

var (
	ErrIsInstalled     = errors.New("already installed")
	ErrIsUpToDate      = errors.New("already up to date")
	ErrInvalidToolName = errors.New("invalid tool name")

	ErrNoAssetFound = "could not find release asset for your platform (%s/%s)"
	ErrToolNotFound = "%s: tool not found in path %s: skipping"
)

type Tool struct {
	Name          string            `json:"name"`
	Repo          string            `json:"repo"`
	Version       string            `json:"version"`
	GoInstallPath string            `json:"go_install_path" yaml:"go_install_path"`
	Requirements  []ToolRequirement `json:"requirements"`
	Assets        map[string]string `json:"assets"`
	InstallType   InstallType       `json:"install_type" yaml:"install_type"`
}

type InstallType string

const (
	Binary InstallType = "binary"
	Go     InstallType = "go"
)

type ToolRequirement struct {
	OS            string                         `json:"os"`
	Specification []ToolRequirementSpecification `json:"specification"`
}

type ToolRequirementSpecification struct {
	Name        string `json:"name"`
	Required    bool   `json:"required"`
	Command     string `json:"command"`
	Instruction string `json:"instruction"`
}

type NucleiData struct {
	IgnoreHash string `json:"ignore-hash"`
	Tools      []Tool `json:"tools"`
}

// IsValidToolName reports whether name is safe to use as a file name in the
// install path: letters, digits, '-' and '_', starting with a letter or a
// digit. Tool names come from the pdtm API or its cache file, and a path
// separator or ".." would reach files outside the install path.
func IsValidToolName(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		case (c == '-' || c == '_') && i > 0:
		default:
			return false
		}
	}
	return true
}

// ValidateToolName returns ErrInvalidToolName for a name that is not safe.
func ValidateToolName(name string) error {
	if !IsValidToolName(name) {
		return fmt.Errorf("%w: %q", ErrInvalidToolName, name)
	}
	return nil
}
