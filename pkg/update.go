package pkg

// Update updates a given tool
func Update(tool Tool, path string) (string, error) {
	if err := Remove(tool); err != nil {
		return "", err
	}
	version, err := Install(tool, path)
	if err != nil {
		return "", err
	}
	return version, nil
}
