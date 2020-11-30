package install

type genericProcess interface {
	Name() (string, error)
	PID() int32
}

type discoverer interface {
	discover() (*discoveryManifest, error)
}
