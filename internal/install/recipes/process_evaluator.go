package recipes

import (
	"context"
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/newrelic/newrelic-cli/internal/install/execution"

	"github.com/shirou/gopsutil/v3/process"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ProcessEvaluatorInterface interface {
	GetOrLoadProcesses(ctx context.Context) []types.GenericProcess
	DetectionStatus(ctx context.Context, r *types.OpenInstallationRecipe, recipeNames []string) execution.RecipeStatusType
	FindProcess(process string) bool
}

type ProcessEvaluator struct {
	processMatchFinder ProcessMatchFinder
	processFetcher     func(context.Context) []types.GenericProcess
	isCached           bool
	cachedProcess      []types.GenericProcess
}

func NewProcessEvaluator() *ProcessEvaluator {
	return newProcessEvaluator(NewRegexProcessMatchFinder(), GetPsUtilCommandLines, true)
}

func newProcessEvaluator(processMatchFinder ProcessMatchFinder, processFetcher func(context.Context) []types.GenericProcess, isCached bool) *ProcessEvaluator {
	return &ProcessEvaluator{
		processMatchFinder: processMatchFinder,
		processFetcher:     processFetcher,
		isCached:           isCached,
		cachedProcess:      nil,
	}
}

func GetPsUtilCommandLines(ctx context.Context) []types.GenericProcess {
	pids, err := process.PidsWithContext(ctx)

	if err != nil {
		log.Errorf("cannot retrieve processes: %s", err)
		return []types.GenericProcess{}
	}

	processes := []types.GenericProcess{}
	for _, pid := range pids {
		var psproc *process.Process
		psproc, err = process.NewProcess(pid)
		if err != nil {
			if err != process.ErrorProcessNotRunning {
				log.Debugf("cannot read pid %d: %s", pid, err)
			}

			continue
		}

		p := NewPSUtilProcess(psproc)
		processes = append(processes, p)
	}

	return processes
}

func (pe *ProcessEvaluator) GetOrLoadProcesses(ctx context.Context) []types.GenericProcess {
	if pe.isCached {
		if pe.cachedProcess == nil {
			pe.cachedProcess = pe.processFetcher(ctx)
		}
		return pe.cachedProcess
	}
	return pe.processFetcher(ctx)
}

func (pe *ProcessEvaluator) DetectionStatus(ctx context.Context, r *types.OpenInstallationRecipe, recipeNames []string) execution.RecipeStatusType {
	if len(r.ProcessMatch) == 0 {
		return execution.RecipeStatusTypes.AVAILABLE
	}

	processes := pe.GetOrLoadProcesses(ctx)
	matches := pe.processMatchFinder.FindMatches(ctx, processes, *r)
	if len(matches) == 0 {
		if slices.Contains(recipeNames, r.Name) {
			log.Errorf("Unsupported (%s): Unable to match any of the following processes:\n", r.DisplayName)
			for _, v := range r.ProcessMatch {
				fmt.Println("-", v)
			}
		}
		log.Tracef("recipe %s is not matching any process", r.Name)
		return execution.RecipeStatusTypes.NULL
	}

	return execution.RecipeStatusTypes.AVAILABLE
}

func (pe *ProcessEvaluator) FindProcess(process string) bool {
	for _, p := range pe.cachedProcess {
		name, _ := p.Name()
		if name == process {
			return true
		}
	}
	return false
}
