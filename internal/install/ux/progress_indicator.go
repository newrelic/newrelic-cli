package ux

type ProgressIndicator interface {
	Fail()
	Success()
	Start(msg string)
	Stop()
}
