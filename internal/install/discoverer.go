package install

import "context"

type discoverer interface {
	discover(context.Context) (*discoveryManifest, error)
}

type discoveryManifest struct {
	Hostname        string           `json:"hostname"`
	KernelArch      string           `json:"kernelArch"`
	KernelVersion   string           `json:"kernelVersion"`
	OS              string           `json:"os"`
	Platform        string           `json:"platform"`
	PlatformFamily  string           `json:"platformFamily"`
	PlatformVersion string           `json:"platformVersion"`
	Processes       []genericProcess `json:"processes"`
}

type genericProcess interface {
	Name() (string, error)
	PID() int32
}

func (d *discoveryManifest) AddProcess(p genericProcess) {
	d.Processes = append(d.Processes, p)
}
