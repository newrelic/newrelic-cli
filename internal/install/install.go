package install

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-task/task/v3"
	taskargs "github.com/go-task/task/v3/args"
	"github.com/go-task/task/v3/taskfile"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func install(configFile string) error {
	// Execute the discovery process.
	log.Debug("Running discovery...")
	var d discoverer = new(mockDiscoverer)
	manifest, err := d.discover()
	if err != nil {
		return err
	}

	// Retrieve the relevant recipes.
	log.Debug("Retrieving recipes...")
	var f recipeFetcher = new(yamlRecipeFetcher)
	recipes, err := f.fetch(configFile, manifest)
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
	}

	if err = e.Setup(); err != nil {
		return err
	}

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
