package output

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
)

// New creates a new outputter with the specific config options
func New(opts ...ConfigOption) (*Output, error) {
	config := Output{
		format:      FormatJSON,
		prettyPrint: true,
	}

	// Loop through config options
	for _, fn := range opts {
		if nil != fn {
			if err := fn(&config); err != nil {
				return nil, err
			}
		}
	}

	switch config.format {
	case FormatJSON:
		config.jsonFormatter = prettyjson.NewFormatter()
		if !config.prettyPrint {
			config.jsonFormatter.DisabledColor = true
			config.jsonFormatter.Indent = 0
			config.jsonFormatter.Newline = ""
		}
	default:
		return nil, fmt.Errorf("unsupported output format %#v", config.format)
	}

	return &config, nil
}

type ConfigOption func(*Output) error

func ConfigFormat(format Format) ConfigOption {
	return func(cfg *Output) error {
		cfg.format = format
		return nil
	}
}

func ConfigPrettyPrint(pretty bool) ConfigOption {
	return func(cfg *Output) error {
		cfg.prettyPrint = pretty
		return nil
	}
}
