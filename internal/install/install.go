package install

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-task/task/v3"
	taskargs "github.com/go-task/task/v3/args"
	"github.com/go-task/task/v3/taskfile"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func install() error {
	// Execute the discovery process.
	log.Debug("Running discovery...")
	d := new(mockDiscoverer)
	manifest, err := d.discover()
	if err != nil {
		return err
	}

	// Retrieve the relevant recipes.
	log.Debug("Retrieving recipes...")
	f := new(yamlRecipeFetcher)
	recipes, err := f.fetch(manifest)
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

	file.Write(out)
	
	e := task.Executor{
		Entrypoint: file.Name(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := e.Setup(); err != nil {
		return err
	}

	var tf taskfile.Taskfile
	err = yaml.Unmarshal(out, &tf)
	if err != nil {
		return err
	}

	taskKeys := getTaskKeys(map[string]*taskfile.Task(tf.Tasks))

	calls, globals := taskargs.ParseV3(taskKeys...)
	e.Taskfile.Vars.Merge(globals)

	if err := e.Run(getSignalContext(), calls...); err != nil {
		return err
	}

	return nil
}

func getTaskKeys(m map[string]*taskfile.Task) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func getSignalContext() context.Context {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-ch
		log.Warnf("signal received: %s", sig)
		cancel()
	}()
	return ctx
}