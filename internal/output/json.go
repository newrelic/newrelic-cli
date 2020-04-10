package output

import (
	"bytes"
	"errors"
	"fmt"
)

// jsonSetPrettyPrint toggles the pretty printing within
// the json formatter
func (o *Output) jsonSetPrettyPrint(pretty bool) error {
	if o == nil || o.jsonFormatter == nil {
		return errors.New("invalid output formatter")
	}

	if pretty {
		o.jsonFormatter.DisabledColor = false
		o.jsonFormatter.Indent = 2
		o.jsonFormatter.Newline = "\n"
	} else {
		o.jsonFormatter.DisabledColor = true
		o.jsonFormatter.Indent = 0
		o.jsonFormatter.Newline = ""
	}

	return nil
}

// JSON prints out data as JSON
func (o *Output) json(data interface{}) error {
	var (
		formatted []byte
		err       error
	)

	// Early quit on no data
	if data == nil {
		return nil
	}

	if o == nil || o.jsonFormatter == nil {
		return errors.New("invalid output formatter")
	}

	// ensure right printing config
	o.jsonSetPrettyPrint(o.prettyPrint)

	// Let's see what they sent us
	switch d := data.(type) {
	case *bytes.Buffer:
		formatted, err = o.jsonFormatter.Format(d.Bytes())
	case []byte:
		formatted, err = o.jsonFormatter.Format(d)
	default:
		formatted, err = o.jsonFormatter.Marshal(d)
	}

	if err != nil {
		return err
	}

	fmt.Println(bytes.NewBuffer(formatted).String())

	return nil
}
