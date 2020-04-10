package output

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

// yaml prints out data as yaml
func (o *Output) yaml(data interface{}) error {
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

	// Let's see what they sent us
	switch d := data.(type) {
	//case *bytes.Buffer:
	//	formatted, err = o.jsonFormatter.Format(d.Bytes())
	//case []byte:
	//	formatted, err = o.jsonFormatter.Format(d)
	default:
		formatted, err = yaml.Marshal(d)
	}

	if err != nil {
		return err
	}

	fmt.Println(string(formatted))

	return nil
}
