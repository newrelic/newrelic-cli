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

var entitiesDescribeTags = &cobra.Command{
	Use:   "describe-tags",
	Short: "Describe the tags for a given entity",
	Long: `Describe the tags for a given entity

The describe-tags command returns JSON output of the tags for the requested
entity.
`,
	Example: "newrelic entities describe-tags --guid <guid>",
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

var entitiesDeleteTags = &cobra.Command{
	Use:   "delete-tags",
	Short: "Delete the given tag:value pairs from the given entitiy",
	Long: `Delete the given tag:value pairs from the given entitiy

The delete-tags command performs a delete operation on the entity for the
specified tag keys.
`,
	Example: "newrelic entities delete-tags --guid <guid> --tag tag1 --tag tag2 --tag tag3,tag4",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			err := nrClient.Entities.DeleteTags(entityGUID, entityTags)
			if err != nil {
				log.Fatal(err)
			}
		})
	},
}

var entitiesDeleteTagValues = &cobra.Command{
	Use:   "delete-tag-values",
	Short: "Delete the given tag/value pairs from the given entitiy",
	Long: `Delete the given tag/value pairs from the given entitiy

The delete-tag-values command performs a delete operation on the entity for the
specified tag:value pairs.
`,
	Example: "newrelic entities delete-tag-values --guid <guid> --tag tag1:value1",
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

var entitiesCreateTags = &cobra.Command{
	Use:   "create-tags",
	Short: "Create tag:value pairs for the given entitiy",
	Long: `Create tag:value pairs for the given entitiy

The create-tags command adds tag:value pairs for the given entity.
`,
	Example: "newrelic entities create-tags --guid <guid> --tag tag1:value1",
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

var entitiesReplaceTags = &cobra.Command{
	Use:   "replace-tags",
	Short: "Replace tag:value pairs for the given entitiy",
	Long: `Repaces tag:value pairs for the given entitiy

The replace-tags command replaces any existing tag:value pairs with those
provided for the given entity.
`,
	Example: "newrelic entities replace-tags --guid <guid> --tag tag1:value1",
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

	Command.AddCommand(entitiesDescribeTags)
	entitiesDescribeTags.Flags().StringVarP(&entityGUID, "guid", "g", "", "entity GUID to describe")
	err = entitiesDescribeTags.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	Command.AddCommand(entitiesDeleteTags)
	entitiesDeleteTags.Flags().StringVarP(&entityGUID, "guid", "g", "", "entity GUID to delete tags on")
	entitiesDeleteTags.Flags().StringSliceVarP(&entityTags, "tag", "t", []string{}, "tag names to delete from the entity")
	err = entitiesDeleteTags.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	err = entitiesDeleteTags.MarkFlagRequired("tag")
	if err != nil {
		log.Error(err)
	}

	Command.AddCommand(entitiesDeleteTagValues)
	entitiesDeleteTagValues.Flags().StringVarP(&entityGUID, "guid", "g", "", "entity GUID to delete tag values on")
	entitiesDeleteTagValues.Flags().StringSliceVarP(&entityValues, "value", "v", []string{}, "key:value tags to delete from the entity")
	err = entitiesDeleteTagValues.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	err = entitiesDeleteTagValues.MarkFlagRequired("value")
	if err != nil {
		log.Error(err)
	}

	Command.AddCommand(entitiesCreateTags)
	entitiesCreateTags.Flags().StringVarP(&entityGUID, "guid", "g", "", "entity GUID to create tag values on")
	entitiesCreateTags.Flags().StringSliceVarP(&entityTags, "tag", "t", []string{}, "tag names to add to the entity")
	err = entitiesCreateTags.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	err = entitiesCreateTags.MarkFlagRequired("tag")
	if err != nil {
		log.Error(err)
	}

	Command.AddCommand(entitiesReplaceTags)
	entitiesReplaceTags.Flags().StringVarP(&entityGUID, "guid", "g", "", "entity GUID to delete tag values on")
	entitiesReplaceTags.Flags().StringSliceVarP(&entityTags, "tag", "t", []string{}, "tag names to replace on the entity")
	err = entitiesReplaceTags.MarkFlagRequired("guid")
	if err != nil {
		log.Error(err)
	}

	err = entitiesReplaceTags.MarkFlagRequired("tag")
	if err != nil {
		log.Error(err)
	}
}
