package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// Command represents the decode command.
var cmdDecode = &cobra.Command{
	Use:     "decode",
	Short:   "Decodes NR1 URL Strings",
	Example: `newrelic decode https://one.newrelic.com/launcher/nr1-core.home?pane=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5ob21lLXNjcmVlbiJ9&platform[accountId]=1`,

	Run: func(cmd *cobra.Command, args []string) {
		relicString := strings.Join(args, " ")
		u, _ := url.Parse(relicString)
		relicString = u.Query().Get("pane")
		decoded, err := base64.StdEncoding.DecodeString(relicString)

		if err != nil {
			fmt.Println("decode error:", err)
		}

		decodedURL := string(decoded)
		fmt.Println("\n**** Decoded URL **** ")
		fmt.Println(decodedURL)

		entityID := strings.Contains(decodedURL, "entityId")

		if entityID {
			var regex = regexp.MustCompile(`[a-zA-Z0-9\+]*={0,3}`)

			entityIDString := gjson.Get(decodedURL, "entityId")
			decodedURL = entityIDString.String()

			if len(regex.FindStringIndex(decodedURL)) > 0 {
				for (len(decodedURL) % 4) != 0 {
					decodedURL += "="
				}
			}

			decoded, err := base64.StdEncoding.DecodeString(decodedURL)
			if err != nil {
				fmt.Println("decode error:", err)
			}

			decodedURL := string(decoded)
			fmt.Println("\n**** Decoded Entity ID **** ")
			fmt.Println(decodedURL)
		}

	},
}

func init() {
	Command.AddCommand(cmdDecode)
}
