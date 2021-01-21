package types

// DiscoveryManifest contains the discovered information about the host.
type DiscoveryManifest struct {
	Hostname        string            `json:"hostname"`
	KernelArch      string            `json:"kernelArch"`
	KernelVersion   string            `json:"kernelVersion"`
	OS              string            `json:"os"`
	Platform        string            `json:"platform"`
	PlatformFamily  string            `json:"platformFamily"`
	PlatformVersion string            `json:"platformVersion"`
	Processes       []ProcessInfoWrap `json:"processes"`
}

// GenericProcess is an abstracted representation of a process.
type GenericProcess interface {
	Name() (string, error)
	Cmdline() (string, error)
	PID() int32
}

type ProcessInfoWrap struct {
	Info            string
	Process         GenericProcess
	MatchingPattern string
}

// AddProcess adds a discovered process to the underlying manifest.
func (d *DiscoveryManifest) AddProcess(p ProcessInfoWrap) {
	d.Processes = append(d.Processes, p)
}
