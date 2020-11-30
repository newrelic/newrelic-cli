package install

import (
	"runtime"
)

type discoveryManifest struct {
	processes         []genericProcess
	platform          string
	arch              string
	targetEnvironment string
}

func newDiscoveryManifest() *discoveryManifest {
	d := discoveryManifest{
		platform:          runtime.GOOS,
		targetEnvironment: "vm",
	}

	return &d
}

func (d *discoveryManifest) ToRecommendationsInput() (*recommendationsInput, error) {
	c := recommendationsInput{
		Variant: variantInput{
			OS:                d.platform,
			Arch:              d.arch,
			TargetEnvironment: d.targetEnvironment,
		},
	}

	for _, process := range d.processes {
		n, err := process.Name()
		if err != nil {
			return nil, err
		}

		p := processDetailInput{
			Name: n,
		}
		c.ProcessDetails = append(c.ProcessDetails, p)
	}

	return &c, nil
}

func (d *discoveryManifest) AddProcess(p genericProcess) {
	d.processes = append(d.processes, p)
}
