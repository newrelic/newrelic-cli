package entities

import (
	"errors"
	"fmt"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

var (
	entityTag  string
	entityTags []string
)

var cmdTags = &cobra.Command{
	Use:   "tags",
	Short: "Manage tags on New Relic entities",
	Long: `Manage entity tags

The tag command allows users to manage the tags applied on the requested
entity. Use --help for more information.
`,
	Example: "newrelic entity tags get --guid <guid>",
}

var cmdTagsGet = &cobra.Command{
	Use:   "get",
	Short: "Get the tags for a given entity",
	Long: `Get the tags for a given entity

The get command returns JSON output of the tags for the requested entity.
`,
	Example: "newrelic entity tags get --guid <entityGUID>",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			tags, err := nrClient.Entities.ListTags(entityGUID)
			if err != nil {
				log.Fatal(err)
			}

			json, err := prettyjson.Marshal(tags)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(json))
		})
	},
}

var cmdTagsDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete the given tag:value pairs from the given entity",
	Long: `Delete the given tag:value pairs from the given entity

The delete command deletes all tags on the given entity 
that match the specified keys.
`,
	Example: "newrelic entity tags delete --guid <entityGUID> --tag tag1 --tag tag2 --tag tag3,tag4",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			err := nrClient.Entities.DeleteTags(entityGUID, entityTags)
			if err != nil {
				log.Fatal(err)
			}
		})
	},
}

var cmdTagsDeleteValues = &cobra.Command{
	Use:   "delete-values",
	Short: "Delete the given tag/value pairs from the given entity",
	Long: `Delete the given tag/value pairs from the given entity

The delete-values command deletes the specified tag:value pairs on a given entity.
`,
	Example: "newrelic entity tags delete-values --guid <guid> --tag tag1:value1",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			tagValues, err := assembleTagValues(entityValues)
			if err != nil {
				log.Fatal(err)
			}

			err = nrClient.Entities.DeleteTagValues(entityGUID, tagValues)
			if err != nil {
				log.Fatal(err)
			}
		})
	},
}

var cmdTagsCreate = &cobra.Command{
	Use:   "create",
	Short: "Create tag:value pairs for the given entity",
	Long: `Create tag:value pairs for the given entity

The create command adds tag:value pairs to the given entity.
`,
	Example: "newrelic entity tags create --guid <entityGUID> --tag tag1:value1",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			tags, err := assembleTags(entityTags)
			if err != nil {
				log.Fatal(err)
			}

			err = nrClient.Entities.AddTags(entityGUID, tags)
			if err != nil {
				log.Fatal(err)
			}
		})
	},
}

var cmdTagsReplace = &cobra.Command{
	Use:   "replace",
	Short: "Replace tag:value pairs for the given entity",
	Long: `Replace tag:value pairs for the given entity

The replace command replaces any existing tag:value pairs with those
provided for the given entity.
`,
	Example: "newrelic entity tags replace --guid <entityGUID> --tag tag1:value1",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			tags, err := assembleTags(entityTags)
			if err != nil {
				log.Fatal(err)
			}

			err = nrClient.Entities.ReplaceTags(entityGUID, tags)
			if err != nil {
				log.Fatal(err)
			}
		})
	},
}

func assembleTags(tags []string) ([]entities.Tag, error) {
	var t []entities.Tag

	tagBuilder := make(map[string][]string)

	for _, x := range tags {
		if !strings.Contains(x, ":") {
			return []entities.Tag{}, errors.New("tags must be specified as colon separated key:value pairs")
		}

		v := strings.SplitN(x, ":", 2)

		tagBuilder[v[0]] = append(tagBuilder[v[0]], v[1])
	}

	for k, v := range tagBuilder {
		tag := entities.Tag{
			Key:    k,
			Values: v,
		}

		t = append(t, tag)
	}

	return t, nil
}

func assembleTagValues(values []string) ([]entities.TagValue, error) {
	var tagValues []entities.TagValue

	for _, x := range values {
		tv, err := assembleTagValue(x)

		if err != nil {
			return []entities.TagValue{}, err
		}

		tagValues = append(tagValues, tv)
	}

	return tagValues, nil
}

func assembleTagValue(tagValueString string) (entities.TagValue, error) {
	tagFormatError := errors.New("tag values must be specified as colon separated key:value pairs")

	if !strings.Contains(tagValueString, ":") {
		return entities.TagValue{}, tagFormatError
	}

	v := strings.SplitN(tagValueString, ":", 2)

	// Handle incomplete tag where the value portion is empty
	if v[1] == "" {
		return entities.TagValue{}, tagFormatError
	}

	tv := entities.TagValue{
		Key:   v[0],
		Value: v[1],
	}

	return tv, nil
}

func init() {
	var err error

	Command.AddCommand(cmdTags)

	cmdTags.AddCommand(cmdTagsGet)
	cmdTagsGet.Flags().StringVarP(&entityGUID, "guid", "g", "", "the entity GUID to retrieve tags for")
	err = cmdTagsGet.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	cmdTags.AddCommand(cmdTagsDelete)
	cmdTagsDelete.Flags().StringVarP(&entityGUID, "guid", "g", "", "the entity GUID to delete tags on")
	cmdTagsDelete.Flags().StringSliceVarP(&entityTags, "tag", "t", []string{}, "the tag keys to delete from the entity")
	err = cmdTagsDelete.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	err = cmdTagsDelete.MarkFlagRequired("tag")
	if err != nil {
		log.Error(err)
	}

	cmdTags.AddCommand(cmdTagsDeleteValues)
	cmdTagsDeleteValues.Flags().StringVarP(&entityGUID, "guid", "g", "", "the entity GUID to delete tag values on")
	cmdTagsDeleteValues.Flags().StringSliceVarP(&entityValues, "value", "v", []string{}, "the tag key:value pairs to delete from the entity")
	err = cmdTagsDeleteValues.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	err = cmdTagsDeleteValues.MarkFlagRequired("value")
	if err != nil {
		log.Error(err)
	}

	cmdTags.AddCommand(cmdTagsCreate)
	cmdTagsCreate.Flags().StringVarP(&entityGUID, "guid", "g", "", "the entity GUID to create tag values on")
	cmdTagsCreate.Flags().StringSliceVarP(&entityTags, "tag", "t", []string{}, "the tag names to add to the entity")
	err = cmdTagsCreate.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	err = cmdTagsCreate.MarkFlagRequired("tag")
	if err != nil {
		log.Error(err)
	}

	cmdTags.AddCommand(cmdTagsReplace)
	cmdTagsReplace.Flags().StringVarP(&entityGUID, "guid", "g", "", "the entity GUID to replace tag values on")
	cmdTagsReplace.Flags().StringSliceVarP(&entityTags, "tag", "t", []string{}, "the tag names to replace on the entity")
	err = cmdTagsReplace.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	err = cmdTagsReplace.MarkFlagRequired("tag")
	if err != nil {
		log.Error(err)
	}
}
