package discovery

import (
	"github.com/shirou/gopsutil/process"
)

type PSUtilProcess process.Process

func (p PSUtilProcess) Name() (string, error) {
	pp := process.Process(p)
	n, err := pp.Name()
	if err != nil {
		return "", err
	}

	return n, nil
}

func (p PSUtilProcess) PID() int32 {
	return process.Process(p).Pid
}
