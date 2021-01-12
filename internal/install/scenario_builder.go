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
	Basic      TestScenario = "BASIC"
	LogMatches TestScenario = "LOG_MATCHES"
	Fail       TestScenario = "FAIL"
)

var (
	TestScenarios = []TestScenario{
		Basic,
		LogMatches,
		Fail,
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
	}

	return nil
}

func (b *ScenarioBuilder) Basic() *RecipeInstaller {

	// mock implementations
	rf := setupRecipeFetcher()
	ers := []execution.StatusReporter{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
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
		statusReporters:   ers,
		prompter:          p,
		progressIndicator: s,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func (b *ScenarioBuilder) Fail() *RecipeInstaller {

	// mock implementations
	rf := setupRecipeFetcher()
	ers := []execution.StatusReporter{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
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
		statusReporters:   ers,
		prompter:          p,
		progressIndicator: s,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func (b *ScenarioBuilder) LogMatches() *RecipeInstaller {

	// mock implementations
	rf := setupRecipeFetcher()
	ers := []execution.StatusReporter{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
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
		statusReporters:   ers,
		prompter:          p,
		progressIndicator: s,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func setupRecipeFetcher() recipes.RecipeFetcher {
	f := recipes.NewMockRecipeFetcher()
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           "Infrastructure Agent Installer",
			ValidationNRQL: "test NRQL",
			File: `
---
name: Infrastructure Agent Installer
install:
  version: "3"
  tasks:
    default:
`,
		},
		{
			Name:           "Logs integration",
			ValidationNRQL: "test NRQL",
			File: `
---
name: Logs integration
install:
  version: "3"
  tasks:
    default:
`,
		},
	}
	f.FetchRecommendationsVal = []types.Recipe{
		{
			Name:           "Recommended recipe",
			ValidationNRQL: "test NRQL",
			File: `
---
name: Recommended recipe
install:
  version: "3"
  tasks:
    default:
`,
		},
	}

	return f
}
