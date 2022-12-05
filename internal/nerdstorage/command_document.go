package nerdstorage

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nerdstorage"
)

var cmdDocument = &cobra.Command{
	Use:     "document",
	Short:   "Read, write, and delete NerdStorage documents.",
	Example: "newrelic nerdstorage document --help",
	Long:    `Read write, and delete NerdStorage documents`,
}

var cmdDocumentGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a NerdStorage document.",
	Long: `Retrieve a NerdStorage document

Retrieve a NerdStorage document.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.
`,
	Example: `
  # Account scope
  newrelic nerdstorage document get --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol --documentId myDoc

  # Entity scope
  newrelic nerdstorage document get --scope ENTITY --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --entityId MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1  --collection myCol --documentId myDoc

  # User scope
  newrelic nerdstorage document get --scope USER --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --collection myCol --documentId myDoc
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var document interface{}
		var err error

		input := nerdstorage.GetDocumentInput{
			PackageID:  packageID,
			Collection: collection,
			DocumentID: documentID,
		}

		switch strings.ToLower(scope) {
		case "account":
			accountID := configAPI.RequireActiveProfileAccountID()
			document, err = client.NRClient.NerdStorage.GetDocumentWithAccountScopeWithContext(utils.SignalCtx, accountID, input)
		case "entity":
			document, err = client.NRClient.NerdStorage.GetDocumentWithEntityScopeWithContext(utils.SignalCtx, entityGUID, input)
		case "user":
			document, err = client.NRClient.NerdStorage.GetDocumentWithUserScopeWithContext(utils.SignalCtx, input)
		default:
			log.Fatal("scope must be one of ACCOUNT, ENTITY, or USER")
		}
		if err != nil {
			log.Fatal(err)
		}

		utils.LogIfFatal(output.Print(document))
		log.Info("success")
	},
}

var cmdDocumentWrite = &cobra.Command{
	Use:   "write",
	Short: "Write a NerdStorage document.",
	Long: `Write a NerdStorage document

Write a NerdStorage document.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.
`,
	Example: `
  # Account scope
  newrelic nerdstorage document write --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol --documentId myDoc --document '{"field": "myValue"}'

  # Entity scope
  newrelic nerdstorage document write --scope ENTITY --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --entityId MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1 --collection myCol --documentId myDoc --document '{"field": "myValue"}'

  # User scope
  newrelic nerdstorage document write --scope USER --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --collection myCol --documentId myDoc --document '{"field": "myValue"}'
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var unmarshaled map[string]interface{}
		err := json.Unmarshal([]byte(document), &unmarshaled)
		if err != nil {
			log.Fatalf("error parsing provided document: %s", err)
		}

		input := nerdstorage.WriteDocumentInput{
			PackageID:  packageID,
			Collection: collection,
			DocumentID: documentID,
			Document:   unmarshaled,
		}

		switch strings.ToLower(scope) {
		case "account":
			accountID := configAPI.RequireActiveProfileAccountID()
			_, err = client.NRClient.NerdStorage.WriteDocumentWithAccountScopeWithContext(utils.SignalCtx, accountID, input)
		case "entity":
			_, err = client.NRClient.NerdStorage.WriteDocumentWithEntityScopeWithContext(utils.SignalCtx, entityGUID, input)
		case "user":
			_, err = client.NRClient.NerdStorage.WriteDocumentWithUserScopeWithContext(utils.SignalCtx, input)
		default:
			log.Fatal("scope must be one of ACCOUNT, ENTITY, or USER")
		}
		if err != nil {
			log.Fatal(err)
		}

		log.Info("success")
	},
}

var cmdDocumentDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a NerdStorage document.",
	Long: `Delete a NerdStorage document

Delete a NerdStorage document.  Valid scopes are ACCOUNT, ENTITY, and USER.
ACCOUNT scope requires a valid account ID and ENTITY scope requires a valid entity
GUID.  A valid Nerdpack package ID is required.
`,
	Example: `
  # Account scope
  newrelic nerdstorage document delete --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol --documentId myDoc

  # Entity scope
  newrelic nerdstorage document delete --scope ENTITY --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --entityId MjUyMDUyOHxFUE18QVBQTElDQVRJT058MjE1MDM3Nzk1 --collection myCol --documentId myDoc

  # User scope
  newrelic nerdstorage document delete --scope USER --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --collection myCol --documentId myDoc
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		input := nerdstorage.DeleteDocumentInput{
			PackageID:  packageID,
			Collection: collection,
			DocumentID: documentID,
		}

		switch strings.ToLower(scope) {
		case "account":
			accountID := configAPI.RequireActiveProfileAccountID()
			_, err = client.NRClient.NerdStorage.DeleteDocumentWithAccountScopeWithContext(utils.SignalCtx, accountID, input)
		case "entity":
			_, err = client.NRClient.NerdStorage.DeleteDocumentWithEntityScopeWithContext(utils.SignalCtx, entityGUID, input)
		case "user":
			_, err = client.NRClient.NerdStorage.DeleteDocumentWithUserScopeWithContext(utils.SignalCtx, input)
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
	Command.AddCommand(cmdDocument)

	cmdDocument.AddCommand(cmdDocumentGet)
	cmdDocumentGet.Flags().StringVarP(&entityGUID, "entityGuid", "e", "", "the entity GUID")
	cmdDocumentGet.Flags().StringVarP(&packageID, "packageId", "p", "", "the external package ID")
	cmdDocumentGet.Flags().StringVarP(&collection, "collection", "c", "", "the collection name to get the document from")
	cmdDocumentGet.Flags().StringVarP(&documentID, "documentId", "d", "", "the document ID")
	cmdDocumentGet.Flags().StringVarP(&scope, "scope", "s", "USER", "the scope to get the document from")

	err := cmdDocumentGet.MarkFlagRequired("packageId")
	utils.LogIfError(err)

	err = cmdDocumentGet.MarkFlagRequired("scope")
	utils.LogIfError(err)

	err = cmdDocumentGet.MarkFlagRequired("collection")
	utils.LogIfError(err)

	err = cmdDocumentGet.MarkFlagRequired("documentId")
	utils.LogIfError(err)

	cmdDocument.AddCommand(cmdDocumentWrite)
	cmdDocumentWrite.Flags().StringVarP(&entityGUID, "entityGuid", "e", "", "the entity GUID")
	cmdDocumentWrite.Flags().StringVarP(&packageID, "packageId", "p", "", "the external package ID")
	cmdDocumentWrite.Flags().StringVarP(&collection, "collection", "c", "", "the collection name to write the document to")
	cmdDocumentWrite.Flags().StringVarP(&documentID, "documentId", "d", "", "the document ID")
	cmdDocumentWrite.Flags().StringVarP(&document, "document", "o", "{}", "the document to be written, in JSON format")
	cmdDocumentWrite.Flags().StringVarP(&scope, "scope", "s", "USER", "the scope to write the document to")

	err = cmdDocumentWrite.MarkFlagRequired("packageId")
	utils.LogIfError(err)

	err = cmdDocumentWrite.MarkFlagRequired("scope")
	utils.LogIfError(err)

	err = cmdDocumentWrite.MarkFlagRequired("document")
	utils.LogIfError(err)

	err = cmdDocumentWrite.MarkFlagRequired("collection")
	utils.LogIfError(err)

	err = cmdDocumentWrite.MarkFlagRequired("documentId")
	utils.LogIfError(err)

	cmdDocument.AddCommand(cmdDocumentDelete)
	cmdDocumentDelete.Flags().StringVarP(&entityGUID, "entityGuid", "e", "", "the entity GUID")
	cmdDocumentDelete.Flags().StringVarP(&packageID, "packageId", "p", "", "the external package ID")
	cmdDocumentDelete.Flags().StringVarP(&collection, "collection", "c", "", "the collection name to delete the document from")
	cmdDocumentDelete.Flags().StringVarP(&documentID, "documentId", "d", "", "the document ID")
	cmdDocumentDelete.Flags().StringVarP(&scope, "scope", "s", "USER", "the scope to delete the document from")

	err = cmdDocumentDelete.MarkFlagRequired("packageId")
	utils.LogIfError(err)

	err = cmdDocumentDelete.MarkFlagRequired("scope")
	utils.LogIfError(err)

	err = cmdDocumentDelete.MarkFlagRequired("collection")
	utils.LogIfError(err)

	err = cmdDocumentDelete.MarkFlagRequired("documentId")
	utils.LogIfError(err)
}
