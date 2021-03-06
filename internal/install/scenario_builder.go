package install

import (
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

type TestScenario string

const (
	Basic        TestScenario = "BASIC"
	LogMatches   TestScenario = "LOG_MATCHES"
	Fail         TestScenario = "FAIL"
	StitchedPath TestScenario = "STITCHED_PATH"
	Canceled     TestScenario = "CANCELED"
)

var (
	TestScenarios = []TestScenario{
		Basic,
		LogMatches,
		Fail,
		StitchedPath,
		Canceled,
	}
	emptyResults = []nrdb.NRDBResult{
		map[string]interface{}{
			"count": 0.0,
		},
	}
	nonEmptyResults = []nrdb.NRDBResult{
		map[string]interface{}{
			"count": 1.0,
		},
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
	case Canceled:
		return b.CanceledInstall()
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
	c := validation.NewMockNRDBClient()
	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)
	v := validation.NewPollingRecipeValidator(c)

	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	s := ux.NewPlainProgress()

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
	c := validation.NewMockNRDBClient()
	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)
	v := validation.NewPollingRecipeValidator(c)

	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewMockFailingRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	pi := ux.NewPlainProgress()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: pi,
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
	c := validation.NewMockNRDBClient()
	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)
	v := validation.NewPollingRecipeValidator(c)
	gff := discovery.NewMockFileFilterer()

	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer(pf)
	re := execution.NewGoTaskRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	pi := ux.NewPlainProgress()

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
		progressIndicator: pi,
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
	c := validation.NewMockNRDBClient()
	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)
	v := validation.NewPollingRecipeValidator(c)

	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	pi := ux.NewPlainProgress()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: pi,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func (b *ScenarioBuilder) CanceledInstall() *RecipeInstaller {
	// mock implementations
	rf := setupRecipeCanceledInstall()
	ers := []execution.StatusSubscriber{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
	statusRollup := execution.NewInstallStatus(ers)
	c := validation.NewMockNRDBClient()
	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)
	v := validation.NewPollingRecipeValidator(c)

	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	pi := ux.NewPlainProgress()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: pi,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func setupRecipeFetcherGuidedInstall() recipes.RecipeFetcher {
	f := recipes.NewMockRecipeFetcher()
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:        "infrastructure-agent-installer",
			DisplayName: "Infrastructure Agent",
			PreInstall: types.OpenInstallationPreInstallConfiguration{
				Info: `
This is the Infrastructure Agent Installer preinstall message.
It is made up of a multi line string.
				`,
			},
			PostInstall: types.OpenInstallationPostInstallConfiguration{
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

func setupRecipeCanceledInstall() recipes.RecipeFetcher {
	f := recipes.NewMockRecipeFetcher()
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           "infrastructure-agent-installer",
			DisplayName:    "Infrastructure Agent",
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
			Name:           "test-canceled-installation",
			DisplayName:    "Test Canceled Installation",
			ValidationNRQL: "test NRQL",
			File: `
---
name: test-canceled-installation
displayName: Test Canceled Installation
description: Scenario to test when a user cancels installation via ctl+c

processMatch: []

validationNrql: "SELECT count(*) from SystemSample where hostname like '{{.HOSTNAME}}' FACET entityGuid SINCE 10 minutes ago"

install:
  version: "3"
  silent: true
  tasks:
    default:
      cmds:
        - task: run
    run:
      cmds:
        - |
          echo "sleeping 10 seconds"
          sleep 10
`,
		},
	}

	return f
}
