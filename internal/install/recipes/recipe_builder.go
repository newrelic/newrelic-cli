package recipes

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type RecipeBuilder struct {
	id                  string
	name                string
	requireAtDiscovery  string
	DiscoveryMode       []types.OpenInstallationDiscoveryMode
	goTaskInstallScript string
	processMatches      []string
	targets             []types.OpenInstallationRecipeInstallTarget
	vars                map[string]string
	dependencies        []*BundleRecipe
	dependencyNames     []string
}

func NewRecipeBuilder() *RecipeBuilder {
	return &RecipeBuilder{
		id:   "id1",
		name: "recipe1",
		vars: make(map[string]string),
		DiscoveryMode: []types.OpenInstallationDiscoveryMode{
			types.OpenInstallationDiscoveryModeTypes.GUIDED,
			types.OpenInstallationDiscoveryModeTypes.TARGETED,
		},
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

func (b *RecipeBuilder) WithDiscoveryMode(discoveryMode []types.OpenInstallationDiscoveryMode) *RecipeBuilder {
	b.DiscoveryMode = discoveryMode
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

func (b *RecipeBuilder) Vars(key string, value string) *RecipeBuilder {
	b.vars[key] = value
	return b
}

func (b *RecipeBuilder) InstallGoTaskScript(script string) *RecipeBuilder {
	b.goTaskInstallScript = script
	return b
}

func (b *RecipeBuilder) InstallShell(script string) *RecipeBuilder {
	goTaskWrap := fmt.Sprintf(`
version: '3'
tasks:
  default:
    cmds:
      - |
        %s
`, script)
	return b.InstallGoTaskScript(goTaskWrap)
}

func (b *RecipeBuilder) Dependency(dependency *BundleRecipe) *RecipeBuilder {
	b.dependencies = append(b.dependencies, dependency)
	return b
}

func (b *RecipeBuilder) DependencyName(depName string) *RecipeBuilder {
	b.dependencyNames = append(b.dependencyNames, depName)
	return b
}

func (b *RecipeBuilder) Build() *types.OpenInstallationRecipe {
	r := &types.OpenInstallationRecipe{
		ID:   b.id,
		Name: b.name,
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: b.requireAtDiscovery,
			DiscoveryMode:      b.DiscoveryMode,
		},
		Install: b.goTaskInstallScript,
	}
	for key, value := range b.vars {
		r.SetRecipeVar(key, value)
	}
	r.ProcessMatch = append(r.ProcessMatch, b.processMatches...)
	r.InstallTargets = append(r.InstallTargets, b.targets...)
	for _, dependency := range b.dependencies {
		if !utils.StringInSlice(dependency.Recipe.Name, r.Dependencies) {
			r.Dependencies = append(r.Dependencies, dependency.Recipe.Name)
		}
	}
	for _, depName := range b.dependencyNames {
		if !utils.StringInSlice(depName, r.Dependencies) {
			r.Dependencies = append(r.Dependencies, depName)
		}
	}
	return r
}

func (b *RecipeBuilder) BuildBundleRecipe() *BundleRecipe {
	r := b.Build()
	br := &BundleRecipe{
		Recipe: r,
	}
	br.Dependencies = append(br.Dependencies, b.dependencies...)
	return br
}
