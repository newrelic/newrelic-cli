package ux

type ProgressIndicator interface {
	Canceled(string)
	Fail(string)
	Success(string)
	Start(string)
	Stop()
}
