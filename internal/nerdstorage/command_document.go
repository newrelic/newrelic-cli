package nerdstorage

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hokaccha/go-prettyjson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/nerdstorage"
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
	Example: `newrelic nerdstorage document get --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol --documentId myDoc`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			var document interface{}
			var err error

			input := nerdstorage.GetDocumentInput{
				PackageID:  packageID,
				Collection: collection,
				DocumentID: documentID,
			}

			switch strings.ToLower(scope) {
			case "account":
				if accountID == 0 {
					log.Fatal("account ID is required when using account scope")
				}

				document, err = nrClient.NerdStorage.GetDocumentWithAccountScope(accountID, input)
			case "entity":
				if entityGUID == "" {
					log.Fatal("entity GUID is required when using entity scope")
				}

				document, err = nrClient.NerdStorage.GetDocumentWithEntityScope(entityGUID, input)
			case "user":
				document, err = nrClient.NerdStorage.GetDocumentWithUserScope(input)
			default:
				log.Fatal("scope must be one of ACCOUNT, ENTITY, or USER")
			}
			if err != nil {
				log.Fatal(err)
			}

			json, err := prettyjson.Marshal(document)
			if err != nil {
				log.Fatal(err)
			}

			log.Info("success")
			fmt.Println(string(json))
		})
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
	Example: `newrelic nerdstorage document write --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol --documentId myDoc --document '{"field": "myValue"}'`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
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
				if accountID == 0 {
					log.Fatal("account ID is required when using account scope")
				}

				_, err = nrClient.NerdStorage.WriteDocumentWithAccountScope(accountID, input)
			case "entity":
				if entityGUID == "" {
					log.Fatal("entity GUID is required when using entity scope")
				}

				_, err = nrClient.NerdStorage.WriteDocumentWithEntityScope(entityGUID, input)
			case "user":
				_, err = nrClient.NerdStorage.WriteDocumentWithUserScope(input)
			default:
				log.Fatal("scope must be one of ACCOUNT, ENTITY, or USER")
			}
			if err != nil {
				log.Fatal(err)
			}

			log.Info("success")
		})
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
	Example: `newrelic nerdstorage document delete --scope ACCOUNT --packageId b0dee5a1-e809-4d6f-bd3c-0682cd079612 --accountId 12345678 --collection myCol --documentId myDoc`,
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			var err error

			input := nerdstorage.DeleteDocumentInput{
				PackageID:  packageID,
				Collection: collection,
				DocumentID: documentID,
			}

			switch strings.ToLower(scope) {
			case "account":
				if accountID == 0 {
					log.Fatal("account ID is required when using account scope")
				}

				_, err = nrClient.NerdStorage.DeleteDocumentWithAccountScope(accountID, input)
			case "entity":
				if entityGUID == "" {
					log.Fatal("entity GUID is required when using entity scope")
				}

				_, err = nrClient.NerdStorage.DeleteDocumentWithEntityScope(entityGUID, input)
			case "user":
				_, err = nrClient.NerdStorage.DeleteDocumentWithUserScope(input)
			default:
				log.Fatal("scope must be one of ACCOUNT, ENTITY, or USER")
			}
			if err != nil {
				log.Fatal(err)
			}

			log.Info("success")
		})
	},
}

func init() {
	Command.AddCommand(cmdDocument)

	cmdDocument.AddCommand(cmdDocumentGet)
	cmdDocumentGet.Flags().IntVar(&accountID, "accountId", 0, "the account ID")
	cmdDocumentGet.Flags().StringVar(&entityGUID, "entityGuid", "", "the entity GUID")
	cmdDocumentGet.Flags().StringVar(&packageID, "packageId", "", "the external package ID")
	cmdDocumentGet.Flags().StringVar(&collection, "collection", "", "the collection name to get the document from")
	cmdDocumentGet.Flags().StringVar(&documentID, "documentId", "", "the document ID")
	cmdDocumentGet.Flags().StringVar(&scope, "scope", "USER", "the scope to get the document from")

	err := cmdDocumentGet.MarkFlagRequired("packageId")
	if err != nil {
		log.Error(err)
	}

	err = cmdDocumentGet.MarkFlagRequired("scope")
	if err != nil {
		log.Error(err)
	}

	err = cmdDocumentGet.MarkFlagRequired("collection")
	if err != nil {
		log.Error(err)
	}

	err = cmdDocumentGet.MarkFlagRequired("documentId")
	if err != nil {
		log.Error(err)
	}

	cmdDocument.AddCommand(cmdDocumentWrite)
	cmdDocumentWrite.Flags().IntVar(&accountID, "accountId", 0, "the account ID")
	cmdDocumentWrite.Flags().StringVar(&entityGUID, "entityGuid", "", "the entity GUID")
	cmdDocumentWrite.Flags().StringVar(&packageID, "packageId", "", "the external package ID")
	cmdDocumentWrite.Flags().StringVar(&collection, "collection", "", "the collection name to write the document to")
	cmdDocumentWrite.Flags().StringVar(&documentID, "documentId", "", "the document ID")
	cmdDocumentWrite.Flags().StringVar(&document, "document", "{}", "the document to be written")
	cmdDocumentWrite.Flags().StringVar(&scope, "scope", "USER", "the scope to write the document to")

	err = cmdDocumentWrite.MarkFlagRequired("packageId")
	if err != nil {
		log.Error(err)
	}

	err = cmdDocumentWrite.MarkFlagRequired("scope")
	if err != nil {
		log.Error(err)
	}

	err = cmdDocumentWrite.MarkFlagRequired("document")
	if err != nil {
		log.Error(err)
	}

	err = cmdDocumentWrite.MarkFlagRequired("collection")
	if err != nil {
		log.Error(err)
	}

	err = cmdDocumentWrite.MarkFlagRequired("documentId")
	if err != nil {
		log.Error(err)
	}

	cmdDocument.AddCommand(cmdDocumentDelete)
	cmdDocumentDelete.Flags().IntVar(&accountID, "accountId", 0, "the account ID")
	cmdDocumentDelete.Flags().StringVar(&entityGUID, "entityGuid", "", "the entity GUID")
	cmdDocumentDelete.Flags().StringVar(&packageID, "packageId", "", "the external package ID")
	cmdDocumentDelete.Flags().StringVar(&collection, "collection", "", "the collection name to delete the document from")
	cmdDocumentDelete.Flags().StringVar(&documentID, "documentId", "", "the document ID")
	cmdDocumentDelete.Flags().StringVar(&scope, "scope", "USER", "the scope to delete the document from")

	err = cmdDocumentDelete.MarkFlagRequired("packageId")
	if err != nil {
		log.Error(err)
	}

	err = cmdDocumentDelete.MarkFlagRequired("scope")
	if err != nil {
		log.Error(err)
	}

	err = cmdDocumentDelete.MarkFlagRequired("collection")
	if err != nil {
		log.Error(err)
	}

	err = cmdDocumentDelete.MarkFlagRequired("documentId")
	if err != nil {
		log.Error(err)
	}
}
