package output

import (
	"github.com/hokaccha/go-prettyjson"

	"github.com/newrelic/newrelic-cli/internal/utils"
)

// globalOutput is the package level config of Output
var globalOutput *Output

// Format provides the list of output formats supported
type Format uint

const (
	FormatJSON Format = iota
	FormatYAML
	//FormatCSV
)

// Output is the main ref for the output package
type Output struct {
	format      Format
	prettyPrint bool

	jsonFormatter *prettyjson.Formatter
}

func SetFormat(format Format) (err error) {
	if err = ensureGlobalOutput(); err != nil {
		return err
	}

	globalOutput.format = format

	return nil
}

func SetPrettyPrint(pretty bool) (err error) {
	if err = ensureGlobalOutput(); err != nil {
		return err
	}

	globalOutput.prettyPrint = pretty

	return nil
}

// ensureGlobalOutput is a helper function to make sure that
// we have a global instance of the outputter at all times
func ensureGlobalOutput() (err error) {
	if globalOutput == nil {
		globalOutput, err = New()
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	utils.LogIfFatal(ensureGlobalOutput())
}
