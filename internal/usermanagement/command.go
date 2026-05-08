package usermanagement

import (
	"github.com/spf13/cobra"
)

var (
	authDomainID string
	userID       string
	userEmail    string
	userName     string
	userType     string
	userTimeZone string
	groupID      string
	groupName    string
)

// Command represents the usermanagement command.
var Command = &cobra.Command{
	Use:   "usermanagement",
	Short: "Manage New Relic users, groups, and authentication domains.",
}
