package basic

func init() {
}

func ExportedA() error {
	return nil
}

func ExportedB() error {
	return nil
}

func unexported() error {
	return nil
}
