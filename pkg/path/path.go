package path

func SetENV(path string) error {
	_, err := add(path)
	if err != nil {
		return err
	}
	return nil
}
