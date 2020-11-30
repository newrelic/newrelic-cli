package install

import "context"

type discoverer interface {
	discover(context.Context) (*discoveryManifest, error)
}

type discoveryManifest struct {
	processes       []genericProcess
	os              string
	platform        string
	platformFamily  string
	platformVersion string
	kernelVersion   string
	kernelArch      string
}

type genericProcess interface {
	Name() (string, error)
	PID() int32
}

func (d *discoveryManifest) AddProcess(p genericProcess) {
	d.processes = append(d.processes, p)
}
