package install

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/go-task/task/v3"
	taskargs "github.com/go-task/task/v3/args"
	"github.com/go-task/task/v3/taskfile"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/credentials"
)

func install(configFiles []string) error {
	// Execute the discovery process.
	log.Debug("Running discovery...")
	var d discoverer = new(diagDiscoverer)
	manifest, err := d.discover()
	if err != nil {
		return err
	}

	log.Debugf("manifest: %+v", manifest)

	// Retrieve the relevant recipes.
	log.Debug("Retrieving recipes...")
	var f recipeFetcher = new(yamlRecipeFetcher)
	recipes, err := f.fetch(configFiles, manifest)
	if err != nil {
		return err
	}

	// Execute the recipe steps.
	for _, r := range recipes {
		err := executeRecipeSteps(r)
		if err != nil {
			return err
		}
	}

	return nil
}

var s *spinner.Spinner

func preRun(t *taskfile.Task, x map[string]interface{}) {
	if t.Name() == "default" {
		return
	}
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = fmt.Sprintf("%s... ", t.Name())
	s.FinalMSG = fmt.Sprintf("%s ...completed.\n", t.Name())

	// x["spinner"] = s
	s.Start()
}

func postRun(t *taskfile.Task, x map[string]interface{}) {
	if t.Name() == "default" {
		return
	}
	// x["spinner"].(*spinner.Spinner).Stop()
	s.Stop()
}

func executeRecipeSteps(r recipe) error {
	log.Debugf("Executing recipe %s", r.Name)

	out, err := yaml.Marshal(r.Install)
	if err != nil {
		return err
	}

	// Create a temporary task file.
	file, err := ioutil.TempFile("", r.Name)
	defer os.Remove(file.Name())
	if err != nil {
		return err
	}

	_, err = file.Write(out)
	if err != nil {
		return err
	}

	e := task.Executor{
		Entrypoint: file.Name(),
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		PreRun:     preRun,
		PostRun:    postRun,
	}

	if err = e.Setup(); err != nil {
		return err
	}

	// e.Stderr = nil
	// e.Stdout = nil
	// e.Stdin = nil

	var tf taskfile.Taskfile
	err = yaml.Unmarshal(out, &tf)
	if err != nil {
		return err
	}

	calls, globals := taskargs.ParseV3()
	e.Taskfile.Vars.Merge(globals)

	credentials.WithCredentials(func(c *credentials.Credentials) {
		v := taskfile.Vars{}
		licenseKey := c.Profiles[c.DefaultProfile].LicenseKey
		if licenseKey == "" {
			err = errors.New("license key not found in default profile")
		}

		v.Set("NR_LICENSE_KEY", taskfile.Var{Static: licenseKey})
		e.Taskfile.Vars.Merge(&v)
	})

	if err != nil {
		return err
	}

	if err := e.Run(getSignalContext(), calls...); err != nil {
		return err
	}

	return nil
}

func getSignalContext() context.Context {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-ch
		log.Warnf("signal received: %s", sig)
		cancel()
	}()
	return ctx
}
