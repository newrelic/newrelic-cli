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
}

func NewProcessEvaluator() *ProcessEvaluator {
	return newProcessEvaluator(NewRegexProcessMatchFinder(), GetPsUtilCommandLines)
}

func newProcessEvaluator(processMatchFinder ProcessMatchFinder, processFetcher func(context.Context) []types.GenericProcess) *ProcessEvaluator {
	return &ProcessEvaluator{
		processMatchFinder: processMatchFinder,
		processFetcher:     processFetcher,
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

func (pe *ProcessEvaluator) DetectionStatus(ctx context.Context, r *types.OpenInstallationRecipe) execution.RecipeStatusType {
	processes := pe.processFetcher(ctx)
	matches := pe.processMatchFinder.FindMatches(ctx, processes, *r)
	filtered := len(r.ProcessMatch) > 0 && len(matches) == 0

	if filtered {
		log.Tracef("recipe %s is not matching any process", r.Name)
		return execution.RecipeStatusTypes.NULL
	}

	return execution.RecipeStatusTypes.AVAILABLE
}
