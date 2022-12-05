package nerdstorage

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nerdstorage"
)

var cmdCollection = &cobra.Command{
	Use:     "collection",
	Short:   "Read, write, and delete NerdStorage collections.",
	Example: "newrelic nerdstorage collection --help",
	Long:    `Read write, and delete NerdStorage collections`,
}

var cmdCollectionGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a NerdStorage collection.",
	Long: `Retrieve a NerdStorage collection

Retrieve a NerdStorage collection.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.
`,
	Example: `
  # Account scope
  newrelic nerdstorage collection get --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol

  # Entity scope
  newrelic nerdstorage collection get --scope ENTITY --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --entityId MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1  --collection myCol

  # User scope
  newrelic nerdstorage collection get --scope USER --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --collection myCol
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var resp []interface{}
		var err error

		input := nerdstorage.GetCollectionInput{
			PackageID:  packageID,
			Collection: collection,
		}

		switch strings.ToLower(scope) {
		case "account":
			accountID := configAPI.RequireActiveProfileAccountID()
			resp, err = client.NRClient.NerdStorage.GetCollectionWithAccountScopeWithContext(utils.SignalCtx, accountID, input)
		case "entity":
			resp, err = client.NRClient.NerdStorage.GetCollectionWithEntityScopeWithContext(utils.SignalCtx, entityGUID, input)
		case "user":
			resp, err = client.NRClient.NerdStorage.GetCollectionWithUserScopeWithContext(utils.SignalCtx, input)
		default:
			log.Fatal("scope must be one of ACCOUNT, ENTITY, or USER")
		}
		if err != nil {
			log.Fatal(err)
		}

		utils.LogIfFatal(output.Print(resp))
		log.Info("success")
	},
}

var cmdCollectionDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a NerdStorage collection.",
	Long: `Delete a NerdStorage collection

Delete a NerdStorage collection.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.
`,
	Example: `
  # Account scope
  newrelic nerdstorage collection delete --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol

  # Entity scope
  newrelic nerdstorage collection delete --scope ENTITY --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --entityId MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1  --collection myCol

  # User scope
  newrelic nerdstorage collection delete --scope USER --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --collection myCol
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		input := nerdstorage.DeleteCollectionInput{
			PackageID:  packageID,
			Collection: collection,
		}

		switch strings.ToLower(scope) {
		case "account":
			accountID := configAPI.RequireActiveProfileAccountID()
			_, err = client.NRClient.NerdStorage.DeleteCollectionWithAccountScopeWithContext(utils.SignalCtx, accountID, input)
		case "entity":
			_, err = client.NRClient.NerdStorage.DeleteCollectionWithEntityScopeWithContext(utils.SignalCtx, entityGUID, input)
		case "user":
			_, err = client.NRClient.NerdStorage.DeleteCollectionWithUserScopeWithContext(utils.SignalCtx, input)
		default:
			log.Fatal("scope must be one of ACCOUNT, ENTITY, or USER")
		}
		if err != nil {
			log.Fatal(err)
		}

		log.Info("success")
	},
}

func init() {
	Command.AddCommand(cmdCollection)

	cmdCollection.AddCommand(cmdCollectionGet)
	cmdCollectionGet.Flags().StringVarP(&entityGUID, "entityGuid", "e", "", "the entity GUID")
	cmdCollectionGet.Flags().StringVarP(&packageID, "packageId", "p", "", "the external package ID")
	cmdCollectionGet.Flags().StringVarP(&collection, "collection", "c", "", "the collection name to get the document from")
	cmdCollectionGet.Flags().StringVarP(&scope, "scope", "s", "USER", "the scope to get the document from")

	err := cmdCollectionGet.MarkFlagRequired("packageId")
	utils.LogIfError(err)

	err = cmdCollectionGet.MarkFlagRequired("scope")
	utils.LogIfError(err)

	err = cmdCollectionGet.MarkFlagRequired("collection")
	utils.LogIfError(err)

	cmdCollection.AddCommand(cmdCollectionDelete)
	cmdCollectionDelete.Flags().StringVarP(&entityGUID, "entityGuid", "e", "", "the entity GUID")
	cmdCollectionDelete.Flags().StringVarP(&packageID, "packageId", "", "p", "the external package ID")
	cmdCollectionDelete.Flags().StringVarP(&collection, "collection", "c", "", "the collection name to delete the document from")
	cmdCollectionDelete.Flags().StringVarP(&scope, "scope", "s", "USER", "the scope to delete the document from")

	err = cmdCollectionDelete.MarkFlagRequired("packageId")
	utils.LogIfError(err)

	err = cmdCollectionDelete.MarkFlagRequired("scope")
	utils.LogIfError(err)

	err = cmdCollectionDelete.MarkFlagRequired("collection")
	utils.LogIfError(err)
}
