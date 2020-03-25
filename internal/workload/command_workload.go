package workload

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
)

var (
	accountID           int
	name                string
	entityGUIDs         []string
	entitySearchQueries []string
	scopeAccountIDs     []int
	guid                string
)

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a New Relic One workload.",
	Long: `Create a New Relic One workload

The create command accepts several different arguments for explicit and dynamic
workloads.  Entity GUIDs can be provided for explicit inclusion of entities, or
entity search queries can be provided for dynamic inclusion of entities.  Multiple
queries will be aggregated together with an OR.  An account scope can optionally
be provided to include entities from different sub-accounts that you also have
access to.
`,
	Example: `newrelic workload create --name 'Example workload' --accountId 12345678 --entitySearchQuery "name like 'Example application'"`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			createInput := workloads.CreateInput{
				Name: name,
			}

			if len(entityGUIDs) > 0 {
				createInput.EntityGUIDs = entityGUIDs
			}

			if len(entitySearchQueries) > 0 {
				var queryInputs []workloads.EntitySearchQueryInput
				for _, q := range entitySearchQueries {
					queryInputs = append(queryInputs, workloads.EntitySearchQueryInput{Query: q})
				}
				createInput.EntitySearchQueries = queryInputs
			}

			if len(scopeAccountIDs) > 0 {
				createInput.ScopeAccountsInput = &workloads.ScopeAccountsInput{AccountIDs: scopeAccountIDs}
			}

			_, err := nrClient.Workloads.CreateWorkload(accountID, createInput)
			if err != nil {
				log.Fatal(err)
			}

			log.Info("success")
		})
	},
}
var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a New Relic One workload.",
	Long: `Delete a New Relic One workload

The delete command accepts a workload's entity GUID.
`,
	Example: `newrelic workload delete --guid 'MjUyMDUyOHxBOE28QVBQTElDQVRDT058MjE1MDM3Nzk1'`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			_, err := nrClient.Workloads.DeleteWorkload(guid)
			if err != nil {
				log.Fatal(err)
			}

			log.Info("success")
		})
	},
}

func init() {
	Command.AddCommand(cmdCreate)
	cmdCreate.Flags().IntVarP(&accountID, "accountId", "a", 0, "the New Relic account ID where you want to create the workload")
	cmdCreate.Flags().StringVarP(&name, "name", "n", "", "the name of the workload")
	cmdCreate.Flags().StringSliceVarP(&entityGUIDs, "entityGuid", "g", []string{}, "the list of entity Guids composing the workload")
	cmdCreate.Flags().StringSliceVarP(&entitySearchQueries, "entitySearchQuery", "q", []string{}, "a list of search queries, combined using an OR operator")
	cmdCreate.Flags().IntSliceVarP(&scopeAccountIDs, "scopeAccountIds", "s", []int{}, "accounts that will be used to get entities from")
	err := cmdCreate.MarkFlagRequired("accountId")
	if err != nil {
		log.Error(err)
	}

	err = cmdCreate.MarkFlagRequired("name")
	if err != nil {
		log.Error(err)
	}

	Command.AddCommand(cmdDelete)
	cmdDelete.Flags().StringVarP(&guid, "guid", "g", "", "the GUID of the workload to delete")
	if err != nil {
		log.Error(err)
	}
}
