package output

// Print outputs the data in the expected format
func Print(data interface{}) (err error) {
	if err = ensureGlobalOutput(); err != nil {
		return err
	}

	switch globalOutput.format {
	case FormatJSON:
		err = globalOutput.json(data)
	case FormatYAML:
		err = globalOutput.yaml(data)
	default:
		err = globalOutput.json(data)
	}

	return err
}

// JSON allows you to override the default output method and
// explicitly print JSON to the screen
func JSON(data interface{}) (err error) {
	if err = ensureGlobalOutput(); err != nil {
		return err
	}

	return globalOutput.json(data)
}
