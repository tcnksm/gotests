package basic

func init() {
	// should be ignored
}

func ExportedA() error {
	return nil
}

func ExportedB() error {
	return nil
}

func Unexported() error {
	return nil
}
