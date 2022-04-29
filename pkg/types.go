package pkg

const Organization = "projectdiscovery"

type Tool struct {
	Name    string            `jðson:"name"`
	Repo    string            `json:"repo"`
	Version string            `json:"version"`
	Assets  map[string]string `json:"assets"`
}
