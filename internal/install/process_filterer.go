package install

type processFilterer interface {
	filter([]genericProcess) ([]genericProcess, error)
}
