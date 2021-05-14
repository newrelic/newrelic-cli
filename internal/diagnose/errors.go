package diagnose

import "errors"

type ErrConnection struct {
	innerErr error
}

func NewErrConnection(err error) ErrConnection {
	return ErrConnection{
		innerErr: err,
	}
}

func (e ErrConnection) Error() string {
	return e.innerErr.Error()
}

func (e ErrConnection) Unwrap() error {
	return e.innerErr
}

type ErrValidation struct {
	innerErr error
}

func NewErrValidation(err error) ErrValidation {
	return ErrValidation{
		innerErr: err,
	}
}

func (e ErrValidation) Error() string {
	return e.innerErr.Error()
}

func (e ErrValidation) Unwrap() error {
	return e.innerErr
}

type ErrDiscovery struct {
	innerErr error
}

func NewErrDiscovery(err error) ErrDiscovery {
	return ErrDiscovery{
		innerErr: err,
	}
}

func (e ErrDiscovery) Error() string {
	return e.innerErr.Error()
}

func (e ErrDiscovery) Unwrap() error {
	return e.innerErr
}

type ErrPostEvent struct {
	innerErr error
}

func NewErrPostEvent(err error) ErrPostEvent {
	return ErrPostEvent{
		innerErr: err,
	}
}

func (e ErrPostEvent) Error() string {
	return e.innerErr.Error()
}

func (e ErrPostEvent) Unwrap() error {
	return e.innerErr
}

var ErrLicenseKey = errors.New("")
var ErrInsightsInsertKey = errors.New("")
