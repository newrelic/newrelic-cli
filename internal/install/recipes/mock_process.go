package recipes

type MockProcess struct {
	cmd  string
	name string
	pid  int32
}

func NewMockProcess(cmd string, name string, pid int32) *MockProcess {
	return &MockProcess{
		cmd:  cmd,
		name: name,
		pid:  pid,
	}
}

func (p MockProcess) Name() (string, error) {
	return p.name, nil
}

func (p MockProcess) Cmd() (string, error) {
	return p.cmd, nil
}

func (p MockProcess) PID() int32 {
	return p.pid
}
