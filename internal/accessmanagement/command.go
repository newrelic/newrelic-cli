package accessmanagement

import (
	"github.com/spf13/cobra"
)

var (
	groupID    string
	roleID     string
	grantScope string
	accountID  int
)

// Command represents the accessmanagement command.
var Command = &cobra.Command{
	Use:   "accessmanagement",
	Short: "Manage New Relic access grants, roles, and permissions.",
}
