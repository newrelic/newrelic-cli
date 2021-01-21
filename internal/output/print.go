package output

import (
	"fmt"
	"log"
)

// Print outputs the data in the expected format
func Print(data interface{}) (err error) {
	if err = ensureGlobalOutput(); err != nil {
		return err
	}

	switch globalOutput.format {
	case FormatJSON:
		err = globalOutput.json(data)
	case FormatText:
		err = globalOutput.text(data)
	case FormatYAML:
		err = globalOutput.yaml(data)
	default:
		err = globalOutput.json(data)
	}

	return err
}

// Printf renders output based on the format and data provided
func Printf(format string, a ...interface{}) {
	if err := ensureGlobalOutput(); err != nil {
		log.Fatal(err)
	}

	data := fmt.Sprintf(format, a...)

	if err := globalOutput.text(data); err != nil {
		log.Fatal(err)
	}
}

// JSON allows you to override the default output method and
// explicitly print JSON to the screen
func JSON(data interface{}) {
	if err := ensureGlobalOutput(); err != nil {
		log.Fatal(err)
	}

	if err := globalOutput.json(data); err != nil {
		log.Fatal(err)
	}
}

// Text allows you to override the default output method and
// explicitly print text to the screen
func Text(data interface{}) {
	if err := ensureGlobalOutput(); err != nil {
		log.Fatal(err)
	}

	if err := globalOutput.text(data); err != nil {
		log.Fatal(err)
	}
}

// YAML allows you to override the default output method and
// explicitly print YAML to the screen
func YAML(data interface{}) {
	if err := ensureGlobalOutput(); err != nil {
		log.Fatal(err)
	}

	if err := globalOutput.yaml(data); err != nil {
		log.Fatal(err)
	}
}
