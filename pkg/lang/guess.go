package lang

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/shirou/gopsutil/process"
)

// Errors
var (
	ErrNotDetermined = errors.New("language cannot be determined")
)

func GetLangs(ctx context.Context) map[ID][]*process.Process {
	results := make(map[ID][]*process.Process)

	l := log.NewStdErr(true)

	pids, err := process.PidsWithContext(ctx)
	if err != nil {
		l.Errorf("cannot retrieve processes: %s", err.Error())
		os.Exit(1)
	}

	//l.Debugf("pids: %+v", pids)

	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			l.Warnf("cannot read pid: %s", pid, err.Error())
			continue
		}

		langID, err := GuessLang(ctx, p)
		if err != nil {
			// process already finished
			if errors.Unwrap(err) == process.ErrorProcessNotRunning {
				continue
			}

			//l.Debugf("error guessing pid %d: lang: %s, err: %s", pid, langID, err)
			continue
		}

		l.Infof("success guessing pid %d: lang: %s", pid, langID)

		results[langID] = append(results[langID], p)
	}

	return results
}

// CmdLineGuessFn guess what language a cmdline runs on.
// fixed error vars are returned:
// - ErrNotDetermined: language cannot be determined
//   this error can also wrap process.ErrorProcessNotRunning when process no longer exists
func GuessLang(ctx context.Context, p *process.Process) (id ID, err error) {
	id = Unknown
	err = ErrNotDetermined

	if p == nil {
		return
	}

	n, errN := p.NameWithContext(ctx)
	if errN != nil {
		err = fmt.Errorf("%s: %w", err, errN)
		return
	}

	// for simple Java PoC hardcoding should be fine
	if n == "java" {
		isIntegration, errE := isAnIntegration(ctx, p)
		if errE != nil {
			err = fmt.Errorf("%s: %w", err, errE)
			return
		}

		if isIntegration {
			err = ErrNotDetermined
			return
		}

		id = Java
		err = nil
	}

	return
}

// isAnIntegration returns true whenever the process is a NR integration, to avoid instrumenting
// instrumentation and awful feedback loops.
func isAnIntegration(ctx context.Context, p *process.Process) (is bool, err error) {
	cliN, errC := p.CmdlineSliceWithContext(ctx)
	if errC != nil {
		err = fmt.Errorf("%s: %w", err, errC)
		return
	}

	if len(cliN) <= 0 {
		err = fmt.Errorf("%s: %s", err, "empty cmdline")
		return
	}

	if strings.Contains(cliN[0], "newrelic-integrations") {
		is = true
		return
	}

	return
}
