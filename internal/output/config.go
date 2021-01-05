package output

import (
	"github.com/hokaccha/go-prettyjson"
	"golang.org/x/term"
)

// New creates a new outputter with the specific config options
func New(opts ...ConfigOption) (*Output, error) {
	config := &Output{
		format:        DefaultFormat,
		prettyPrint:   DefaultPretty,
		terminalWidth: DefaultTerminalWidth,
	}

	// Set some defaults
	w, _, err := term.GetSize(0)
	if err == nil {
		config.terminalWidth = w
	}

	// Loop through config options
	for _, fn := range opts {
		if nil != fn {
			if err := fn(config); err != nil {
				return nil, err
			}
		}
	}

	// JSON output
	config.jsonFormatter = prettyjson.NewFormatter()

	return config, nil
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
