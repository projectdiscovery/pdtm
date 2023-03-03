package pkg

import (
	"errors"
)

const Organization = "projectdiscovery"

var (
	ErrIsInstalled = errors.New("already installed")
	ErrIsUpToDate  = errors.New("already up to date")

	ErrNoAssetFound = "could not find release asset for your platform (%s/%s)"
	ErrToolNotFound = "tool %s not found in path %s: skipping"
)

type Tool struct {
	Name    string            `json:"name"`
	Repo    string            `json:"repo"`
	Version string            `json:"version"`
	Assets  map[string]string `json:"assets"`
}
