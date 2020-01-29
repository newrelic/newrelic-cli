package main

import (
	log "github.com/sirupsen/logrus"

	// Main entry point is cmd
	"github.com/newrelic/newrelic-cli/internal/cmd"

	// Commands to import, init will run and register with cmd
	_ "github.com/newrelic/newrelic-cli/internal/entities"
)

var (
	// AppName for this CMD
	AppName = "newrelic-dev"
	// Version of the CLI
	Version = "dev"
)

func main() {
	err := cmd.Execute(AppName, Version)
	if err != nil {
		log.Fatal(err)
	}
}
