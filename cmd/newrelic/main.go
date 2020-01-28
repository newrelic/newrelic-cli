package main

import (
	"github.com/newrelic/newrelic-cli/internal/cmd"

	_ "github.com/newrelic/newrelic-cli/internal/entities"
)

func main() {
	cmd.Execute()
}
