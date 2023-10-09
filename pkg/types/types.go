package types

import "errors"

const Organization = "projectdiscovery"

var (
	ErrIsInstalled = errors.New("already installed")
	ErrIsUpToDate  = errors.New("already up to date")

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
