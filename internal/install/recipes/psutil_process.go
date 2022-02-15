package recipes

import (
	"github.com/shirou/gopsutil/process"
)

type PSUtilProcess struct {
	proc *process.Process
}

func NewPSUtilProcess(p *process.Process) PSUtilProcess {
	return PSUtilProcess{
		proc: p,
	}
}

func (p PSUtilProcess) Name() (string, error) {
	n, err := p.proc.Name()
	if err != nil {
		return "", err
	}

	return n, nil
}

func (p PSUtilProcess) Cmd() (string, error) {
	n, err := p.proc.Cmdline()
	if err != nil {
		return "", err
	}

	return n, nil
}

func (p PSUtilProcess) PID() int32 {
	return p.proc.Pid
}
