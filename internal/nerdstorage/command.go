package nerdstorage

import (
	"github.com/spf13/cobra"
)

var (
	entityGUID string
	packageID  string
	collection string
	documentID string
	document   string
	scope      string
)

// Command represents the nerdstorage command.
var Command = &cobra.Command{
	Use:   "nerdstorage",
	Short: "Read, write, and delete NerdStorage documents and collections.",
}
