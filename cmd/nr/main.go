package main

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/pkg/config"
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

	fmt.Printf("%+v\n", cfg)
}
