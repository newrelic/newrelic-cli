package types

import (
	"errors"
)

// ErrInterrupt represents a context cancellation.
var ErrInterrupt = errors.New("operation canceled")
