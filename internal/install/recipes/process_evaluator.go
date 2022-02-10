package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

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

		p := discovery.NewPSUtilProcess(psproc)
		processes = append(processes, p)
	}

	return processes
}

func (pe *ProcessEvaluator) getOrLoadProcesses(ctx context.Context) []types.GenericProcess {
	if (pe.isCached) {
		if (pe.cachedProcess == nil) {
			pe.cachedProcess = pe.processFetcher(ctx)
		}
		return pe.cachedProcess
	}
	return pe.processFetcher(ctx)
}

func (pe *ProcessEvaluator) DetectionStatus(ctx context.Context, r *types.OpenInstallationRecipe) execution.RecipeStatusType {
	if len(r.ProcessMatch) == 0 {
		return execution.RecipeStatusTypes.AVAILABLE
	}

	processes := pe.getOrLoadProcesses(ctx)
	matches := pe.processMatchFinder.FindMatches(ctx, processes, *r)
	if len(matches) == 0 {
		log.Tracef("recipe %s is not matching any process", r.Name)
		return execution.RecipeStatusTypes.NULL
	}

	return execution.RecipeStatusTypes.AVAILABLE
}
