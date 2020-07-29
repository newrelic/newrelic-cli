package output

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/utils"
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
	utils.LogIfFatal(ensureGlobalOutput())

	data := fmt.Sprintf(format, a...)

	utils.LogIfFatal(globalOutput.text(data))
}

// JSON allows you to override the default output method and
// explicitly print JSON to the screen
func JSON(data interface{}) {
	utils.LogIfFatal(ensureGlobalOutput())
	utils.LogIfFatal(globalOutput.json(data))
}

// Text allows you to override the default output method and
// explicitly print text to the screen
func Text(data interface{}) {
	utils.LogIfFatal(ensureGlobalOutput())
	utils.LogIfFatal(globalOutput.text(data))
}

// YAML allows you to override the default output method and
// explicitly print YAML to the screen
func YAML(data interface{}) {
	utils.LogIfFatal(ensureGlobalOutput())
	utils.LogIfFatal(globalOutput.yaml(data))
}
