package recipes

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeBuilder struct {
	id                 string
	name               string
	requireAtDiscovery string
	processMatches     []string
	targets            []types.OpenInstallationRecipeInstallTarget
}

func NewRecipeBuilder() *RecipeBuilder {
	return &RecipeBuilder{
		id:   "id1",
		name: "recipe1",
	}
}

func (b *RecipeBuilder) ID(id string) *RecipeBuilder {
	b.id = id
	return b
}

func (b *RecipeBuilder) Name(name string) *RecipeBuilder {
	b.name = name
	return b
}

func (b *RecipeBuilder) WithPreInstallScript(script string) *RecipeBuilder {
	b.requireAtDiscovery = script
	return b
}

func (b *RecipeBuilder) ProcessMatch(match string) *RecipeBuilder {
	b.processMatches = append(b.processMatches, match)
	return b
}

func (b *RecipeBuilder) TargetOs(os types.OpenInstallationOperatingSystem) *RecipeBuilder {
	t := types.OpenInstallationRecipeInstallTarget{
		Os: os,
	}
	b.targets = append(b.targets, t)
	return b
}

func (b *RecipeBuilder) TargetOsPlatform(os types.OpenInstallationOperatingSystem, platform types.OpenInstallationPlatform) *RecipeBuilder {
	t := types.OpenInstallationRecipeInstallTarget{
		Os:       os,
		Platform: platform,
	}
	b.targets = append(b.targets, t)
	return b
}

func (b *RecipeBuilder) TargetOsPlatformVersionArch(os types.OpenInstallationOperatingSystem, platformVersion string, arch string) *RecipeBuilder {
	t := types.OpenInstallationRecipeInstallTarget{
		Os:              os,
		PlatformVersion: platformVersion,
		KernelArch:      arch,
	}
	b.targets = append(b.targets, t)
	return b
}

func (b *RecipeBuilder) TargetOsArch(os types.OpenInstallationOperatingSystem, arch string) *RecipeBuilder {
	return b.TargetOsPlatformVersionArch(os, "", arch)
}

func (b *RecipeBuilder) Build() *types.OpenInstallationRecipe {
	r := &types.OpenInstallationRecipe{
		ID:   b.id,
		Name: b.name,
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: b.requireAtDiscovery,
		},
	}
	r.ProcessMatch = append(r.ProcessMatch, b.processMatches...)
	r.InstallTargets = append(r.InstallTargets, b.targets...)
	return r
}
