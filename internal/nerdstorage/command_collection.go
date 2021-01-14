package nerdstorage

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/pkg/nerdstorage"
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
	Run: func(cmd *cobra.Command, args []string) {
		nrClient, err := client.NewClient(configuration.GetActiveProfileName())
		if err != nil {
			log.Fatal(err)
		}

		var resp []interface{}

		input := nerdstorage.GetCollectionInput{
			PackageID:  packageID,
			Collection: collection,
		}

		switch strings.ToLower(scope) {
		case "account":
			resp, err = nrClient.NerdStorage.GetCollectionWithAccountScope(accountID, input)
		case "entity":
			resp, err = nrClient.NerdStorage.GetCollectionWithEntityScope(entityGUID, input)
		case "user":
			resp, err = nrClient.NerdStorage.GetCollectionWithUserScope(input)
		default:
			log.Fatal("scope must be one of ACCOUNT, ENTITY, or USER")
		}
		if err != nil {
			log.Fatal(err)
		}

		if err := output.Print(resp); err != nil {
			log.Fatal(err)
		}

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
	Run: func(cmd *cobra.Command, args []string) {
		nrClient, err := client.NewClient(configuration.GetActiveProfileName())
		if err != nil {
			log.Fatal(err)
		}

		input := nerdstorage.DeleteCollectionInput{
			PackageID:  packageID,
			Collection: collection,
		}

		switch strings.ToLower(scope) {
		case "account":
			_, err = nrClient.NerdStorage.DeleteCollectionWithAccountScope(accountID, input)
		case "entity":
			_, err = nrClient.NerdStorage.DeleteCollectionWithEntityScope(entityGUID, input)
		case "user":
			_, err = nrClient.NerdStorage.DeleteCollectionWithUserScope(input)
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
	cmdCollectionGet.Flags().IntVarP(&accountID, "accountId", "a", 0, "the account ID")
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
	cmdCollectionDelete.Flags().IntVarP(&accountID, "accountId", "a", 0, "the account ID")
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
