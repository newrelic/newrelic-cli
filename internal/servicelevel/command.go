package servicelevel

import (
	"github.com/spf13/cobra"
)

// Command represents the servicelevel command.
var Command = &cobra.Command{
	Use:   "servicelevel",
	Short: "Interact with New Relic servicelevel",
}
