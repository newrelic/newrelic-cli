package install

import (
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
)

type TestScenario string

const (
	Basic        TestScenario = "BASIC"
	LogMatches   TestScenario = "LOG_MATCHES"
	Fail         TestScenario = "FAIL"
	StitchedPath TestScenario = "STITCHED_PATH"
)

var (
	TestScenarios = []TestScenario{
		Basic,
		LogMatches,
		Fail,
		StitchedPath,
	}
)

func TestScenarioValues() []string {
	v := make([]string, len(TestScenarios))
	for i, s := range TestScenarios {
		v[i] = string(s)
	}

	return v
}

type ScenarioBuilder struct {
	installerContext InstallerContext
}

func NewScenarioBuilder(ic InstallerContext) *ScenarioBuilder {
	b := ScenarioBuilder{
		installerContext: ic,
	}

	return &b
}

func (b *ScenarioBuilder) BuildScenario(s TestScenario) *RecipeInstaller {
	switch s {
	case Basic:
		return b.Basic()
	case LogMatches:
		return b.LogMatches()
	case Fail:
		return b.Fail()
	case StitchedPath:
		return b.StitchedPath()
	}

	return nil
}

func (b *ScenarioBuilder) Basic() *RecipeInstaller {

	// mock implementations
	rf := setupRecipeFetcherGuidedInstall()
	ers := []execution.StatusSubscriber{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
	statusRollup := execution.NewInstallStatus(ers)
	v := validation.NewMockRecipeValidator()

	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	s := ux.NewSpinner()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: s,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func (b *ScenarioBuilder) Fail() *RecipeInstaller {

	// mock implementations
	rf := setupRecipeFetcherGuidedInstall()
	ers := []execution.StatusSubscriber{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
	statusRollup := execution.NewInstallStatus(ers)
	v := validation.NewMockRecipeValidator()

	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewMockFailingRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	s := ux.NewSpinner()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: s,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func (b *ScenarioBuilder) LogMatches() *RecipeInstaller {

	// mock implementations
	rf := setupRecipeFetcherGuidedInstall()
	ers := []execution.StatusSubscriber{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
	statusRollup := execution.NewInstallStatus(ers)
	v := validation.NewMockRecipeValidator()
	gff := discovery.NewMockFileFilterer()

	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer(pf)
	re := execution.NewGoTaskRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	s := ux.NewSpinner()

	gff.FilterVal = []types.LogMatch{
		{
			Name: "asdf",
			File: "asdf",
		},
	}

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: s,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func (b *ScenarioBuilder) StitchedPath() *RecipeInstaller {
	// mock implementations
	rf := setupRecipeFetcherStitchedPath()
	ers := []execution.StatusSubscriber{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
	statusRollup := execution.NewInstallStatus(ers)
	v := validation.NewMockRecipeValidator()

	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	s := ux.NewSpinner()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: s,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func setupRecipeFetcherGuidedInstall() recipes.RecipeFetcher {
	f := recipes.NewMockRecipeFetcher()
	f.FetchRecipeVals = []types.Recipe{
		{
			Name: "Infrastructure Agent Installer",
			PreInstall: types.OpenInstallationPreInstallConfiguration{
				Info: `
This is the Infrastructure Agent Installer preinstall message.
It is made up of a multi line string.
				`,
			},
			PostInstall: types.RecipePostInstall{
				Info: `
This is the Infrastructure Agent Installer postinstall message.
It is made up of a multi line string.
				`,
			},
			ValidationNRQL: "test NRQL",
			File: `
---
name: infra-agent
displayName: Infrastructure Agent
install:
  version: "3"
  tasks:
    default:
`,
		},
		{
			Name:           "logs-integration",
			DisplayName:    "Logs integration",
			ValidationNRQL: "test NRQL",
			File: `
---
name: logs-integration
displayName: Logs integration
install:
  version: "3"
  tasks:
    default:
`,
		},
	}
	f.FetchRecommendationsVal = []types.Recipe{
		{
			Name:           "recommended-recipe",
			DisplayName:    "Recommended recipe",
			ValidationNRQL: "test NRQL",
			File: `
---
name: recommended-recipe
displayName: Recommended recipe
install:
  version: "3"
  tasks:
    default:
`,
		},
	}

	return f
}

func setupRecipeFetcherStitchedPath() recipes.RecipeFetcher {
	f := recipes.NewMockRecipeFetcher()
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           "recommended-recipe",
			DisplayName:    "Recommended recipe",
			ValidationNRQL: "test NRQL",
			File: `
---
name: recommended-recipe
displayName: Recommended recipe
install:
  version: "3"
  tasks:
    default:
`,
		},
		{
			Name:           "another-recommended-recipe",
			DisplayName:    "Another Recommended recipe",
			ValidationNRQL: "test NRQL",
			File: `
---
name: another-recommended-recipe
displayName: Another Recommended recipe
install:
  version: "3"
  tasks:
    default:
`,
		},
	}

	return f
}
