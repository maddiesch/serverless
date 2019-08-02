package serverless

// Must takes an error and panics
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
