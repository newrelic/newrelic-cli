package types

// ErrInterrupt represents a context cancellation.
type ErrInterrupt struct{}

func (e *ErrInterrupt) Error() string {
	return "operation cancelled"
}

// NewErrInterrupt creates a new instance of ErrInterrupt.
func NewErrInterrupt() *ErrInterrupt {
	return &ErrInterrupt{}
}
