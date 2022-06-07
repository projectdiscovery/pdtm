package pkg

import "errors"

const Organization = "projectdiscovery"

var (
	ErrIsInstalled = errors.New("already installed")
	ErrIsUpToDate  = errors.New("already up to date")
)

type Tool struct {
	Name    string            `jðson:"name"`
	Repo    string            `json:"repo"`
	Version string            `json:"version"`
	Assets  map[string]string `json:"assets"`
}
