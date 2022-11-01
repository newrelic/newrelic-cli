package yq

// Borrowed from gojq - cli package
// ref: https://github.com/itchyny/gojq/blob/main/cli/marshaler.go

import (
	"io"

	"gopkg.in/yaml.v3"
)

type Marshaler interface {
	Marshal(interface{}, io.Writer) error
}

func YamlFormatter(indent *int) *YamlMarshaler {
	return &YamlMarshaler{indent}
}

type YamlMarshaler struct {
	indent *int
}

func (m *YamlMarshaler) Marshal(v interface{}, w io.Writer) error {
	enc := yaml.NewEncoder(w)
	if i := m.indent; i != nil {
		enc.SetIndent(*i)
	} else {
		enc.SetIndent(2)
	}
	if err := enc.Encode(v); err != nil {
		return err
	}
	return enc.Close()
}
