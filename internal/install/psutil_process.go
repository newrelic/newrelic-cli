package install

import (
	"github.com/shirou/gopsutil/process"
)

type psUtilProcess process.Process

func (p psUtilProcess) Name() (string, error) {
	pp := process.Process(p)
	n, err := pp.Name()
	if err != nil {
		return "", err
	}

	return n, nil
}

func (p psUtilProcess) PID() int32 {
	return process.Process(p).Pid
}
