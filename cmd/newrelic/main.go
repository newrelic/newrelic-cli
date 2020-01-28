package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/cmd"
	"github.com/newrelic/newrelic-cli/pkg/config"

	// Commands
	_ "github.com/newrelic/newrelic-cli/internal/entities"
)

const (
	cmdName = "nr"
	appName = "New Relic CLI"
)

var (
	// Version of this command
	Version = "dev"
)

func main() {
	fmt.Printf("%s version: '%s'\n", appName, Version)

	cfg, err := config.Load("", "debug")
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err == nil {
		log.SetLevel(lvl)
	}

	// Main entry point for the CLI
	cmd.Execute()
}
