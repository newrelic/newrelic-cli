package types

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// DiscoveryManifest contains the discovered information about the host.
type DiscoveryManifest struct {
	Hostname        string           `json:"hostname"`
	KernelArch      string           `json:"kernelArch"`
	KernelVersion   string           `json:"kernelVersion"`
	OS              string           `json:"os"`
	Platform        string           `json:"platform"`
	PlatformFamily  string           `json:"platformFamily"`
	PlatformVersion string           `json:"platformVersion"`
	Processes       []MatchedProcess `json:"processes"`
}

// GenericProcess is an abstracted representation of a process.
type GenericProcess interface {
	Name() (string, error)
	Cmdline() (string, error)
	PID() int32
}

type MatchedProcess struct {
	Command         string `json:"command"`
	Process         GenericProcess
	MatchingPattern string
}

// AddMatchedProcess adds a discovered process to the underlying manifest.
func (d *DiscoveryManifest) AddMatchedProcess(p MatchedProcess) {
	d.Processes = append(d.Processes, p)
}

func (d *DiscoveryManifest) ConstrainRecipes(allRecipes []Recipe) []Recipe {
	var recipes []Recipe

	for _, recipe := range allRecipes {
		if len(recipe.InstallTargets) == 0 {
			log.Warnf("recipe has no InstallTargets: %s", recipe.Name)
		}

		for _, target := range recipe.InstallTargets {
			if target.KernelArch != "" {
				if strings.EqualFold(target.KernelArch, d.KernelArch) {
					recipes = append(recipes, recipe)
					break
				}
			}

			if target.KernelVersion != "" {
				if strings.EqualFold(target.KernelVersion, d.KernelVersion) {
					recipes = append(recipes, recipe)
					break
				}
			}

			if target.Os != "" {
				if strings.EqualFold(string(target.Os), d.OS) {
					recipes = append(recipes, recipe)
					break
				}
			}

			if target.Platform != "" {
				if strings.EqualFold(string(target.Platform), d.Platform) {
					recipes = append(recipes, recipe)
					break
				}
			}

			if target.PlatformFamily != "" {
				if strings.EqualFold(string(target.PlatformFamily), d.PlatformFamily) {
					recipes = append(recipes, recipe)
					break
				}
			}

			if target.PlatformVersion != "" {
				if strings.EqualFold(target.PlatformVersion, d.PlatformVersion) {
					recipes = append(recipes, recipe)
					break
				}
			}
		}
	}

	log.Debugf("%d embedded recipes found for manifest", len(recipes))

	return recipes
}
