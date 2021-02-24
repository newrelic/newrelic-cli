package ux

type ProgressIndicator interface {
	Fail(string)
	Success(string)
	Start(string)
	Stop()
}
