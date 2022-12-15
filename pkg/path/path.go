package path

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
