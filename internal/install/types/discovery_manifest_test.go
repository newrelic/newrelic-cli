// +build unit

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscoveryManifest_ConstrainRecipes(t *testing.T) {

	recipes := []OpenInstallationRecipe{
		{
			Name: "unknown",
			InstallTargets: []OpenInstallationRecipeInstallTarget{
				{
					Os: OpenInstallationOperatingSystemTypes.LINUX,
				},
			},
		},
		{
			Name: "windows-test",
			InstallTargets: []OpenInstallationRecipeInstallTarget{
				{
					Os: OpenInstallationOperatingSystemTypes.WINDOWS,
				},
			},
		},
		{
			Name: "suse-sample",
			InstallTargets: []OpenInstallationRecipeInstallTarget{
				{
					Os:       OpenInstallationOperatingSystemTypes.LINUX,
					Platform: OpenInstallationPlatformTypes.SUSE,
					Type:     OpenInstallationTargetTypeTypes.HOST,
				},
			},
		},
		{
			Name: "amd64-linux",
			InstallTargets: []OpenInstallationRecipeInstallTarget{
				{
					KernelArch: "amd64",
					Os:         OpenInstallationOperatingSystemTypes.LINUX,
				},
			},
		},
		{
			Name: "amd64-linux-1-debian",
			InstallTargets: []OpenInstallationRecipeInstallTarget{
				{
					KernelArch:     "amd64",
					Os:             OpenInstallationOperatingSystemTypes.LINUX,
					KernelVersion:  "1.0",
					PlatformFamily: OpenInstallationPlatformFamilyTypes.DEBIAN,
				},
			},
		},
		{
			Name: "amd64-linux-debian/ubuntu/amazon",
			InstallTargets: []OpenInstallationRecipeInstallTarget{
				{
					KernelArch: "amd64",
					Os:         OpenInstallationOperatingSystemTypes.LINUX,
					Platform:   OpenInstallationPlatformTypes.DEBIAN,
				},
				{
					KernelArch: "amd64",
					Os:         OpenInstallationOperatingSystemTypes.LINUX,
					Platform:   OpenInstallationPlatformTypes.AMAZON,
				},
				{
					KernelArch: "amd64",
					Os:         OpenInstallationOperatingSystemTypes.LINUX,
					Platform:   OpenInstallationPlatformTypes.UBUNTU,
				},
			},
		},
		{
			Name: "1.0",
			InstallTargets: []OpenInstallationRecipeInstallTarget{
				{
					PlatformVersion: "1.0",
				},
			},
		},
	}

	cases := []struct {
		manifest DiscoveryManifest
		results  []string
	}{
		{
			manifest: DiscoveryManifest{
				OS: "linux",
			},
			results: []string{"unknown"},
		},
		{
			manifest: DiscoveryManifest{
				OS: "darwin",
			},
			results: []string{},
		},
		{
			manifest: DiscoveryManifest{
				OS: "windows",
			},
			results: []string{"windows-test"},
		},
		{
			manifest: DiscoveryManifest{
				OS:             "linux",
				KernelArch:     "amd64",
				KernelVersion:  "1.0",
				PlatformFamily: "1.0",
				Platform:       "suse",
			},
			results: []string{"unknown", "suse-sample", "amd64-linux"},
		},
		{
			manifest: DiscoveryManifest{
				KernelArch: "amd64",
				OS:         "linux",
				Platform:   "debian",
			},
			results: []string{"unknown", "amd64-linux", "amd64-linux-debian/ubuntu/amazon"},
		},
		{
			manifest: DiscoveryManifest{
				KernelArch:     "amd64",
				KernelVersion:  "1.0",
				OS:             "linux",
				Platform:       "debian",
				PlatformFamily: "debian",
			},
			results: []string{"unknown", "amd64-linux", "amd64-linux-1-debian", "amd64-linux-debian/ubuntu/amazon"},
		},
	}

	for _, c := range cases {
		r := c.manifest.ConstrainRecipes(recipes)

		names := []string{}

		for _, n := range r {
			names = append(names, n.Name)
		}

		require.Equal(t, c.results, names)
	}

}
