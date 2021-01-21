package output

import (
	"log"
	"strings"

	"github.com/hokaccha/go-prettyjson"
)

// globalOutput is the package level config of Output
var globalOutput *Output

// Format provides the list of output formats supported
type Format uint

const DefaultFormat = FormatJSON
const DefaultPretty = true
const DefaultTerminalWidth = 80

const (
	FormatJSON Format = iota
	FormatText
	FormatYAML
	//FormatCSV
)

var formatStrings = map[Format]string{
	FormatJSON: "JSON",
	FormatText: "Text",
	FormatYAML: "YAML",
}

// Output is the main ref for the output package
type Output struct {
	format        Format
	prettyPrint   bool
	terminalWidth int

	jsonFormatter *prettyjson.Formatter
}

// String returns the string value of the format name
func (f Format) String() string {
	if name, ok := formatStrings[f]; ok {
		return name
	}

	return ""
}

func FormatOptions() string {
	ret := make([]string, 0, len(formatStrings))

	for _, v := range formatStrings {
		ret = append(ret, v)
	}
	return strings.Join(ret, ", ")
}

func ParseFormat(name string) Format {
	for k, v := range formatStrings {
		if strings.EqualFold(name, v) {
			return k
		}
	}

	return DefaultFormat
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
	if err := ensureGlobalOutput(); err != nil {
		log.Fatal(err)
	}
}
