package types

// DiscoveryManifest contains the discovered information about the host.
type DiscoveryManifest struct {
	Hostname        string           `json:"hostname"`
	KernelArch      string           `json:"kernelArch"`
	KernelVersion   string           `json:"kernelVersion"`
	OS              string           `json:"os"`
	Platform        string           `json:"platform"`
	PlatformFamily  string           `json:"platformFamily"`
	PlatformVersion string           `json:"platformVersion"`
	Processes       []GenericProcess `json:"processes"`
}

// GenericProcess is an abstracted representation of a process.
type GenericProcess interface {
	Name() (string, error)
	PID() int32
}

// AddProcess adds a discovered process to the underlying manifest.
func (d *DiscoveryManifest) AddProcess(p GenericProcess) {
	d.Processes = append(d.Processes, p)
}
