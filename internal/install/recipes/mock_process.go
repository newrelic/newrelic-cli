package recipes

type mockProcess struct {
	cmdline string
	name    string
	pid     int32
}

func (p mockProcess) Name() (string, error) {
	return p.name, nil
}

func (p mockProcess) Cmd() (string, error) {
	return p.cmdline, nil
}

func (p mockProcess) PID() int32 {
	return p.pid
}
