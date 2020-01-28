package main

import (
	// Main entry point is cmd
	"github.com/newrelic/newrelic-cli/internal/cmd"

	// Commands to import, init will run and register with cmd
	_ "github.com/newrelic/newrelic-cli/internal/entities"
)

// Version of the CLI
var Version = "dev"

func main() {
	cmd.Execute(Version)
}
