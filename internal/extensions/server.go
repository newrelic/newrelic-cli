package extensions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultListenPort = 8080
)

type Route struct {
	Route   string
	Handler routeHandlerFunc
}

type routeHandlerFunc func(cmdDef *CommandDefinition, command *cobra.Command, args []string) func(http.ResponseWriter, *http.Request)

type CommandableResponse struct {
	Cmd         string         `json:"cmd"`
	Args        []string       `json:"args"`
	Flags       []*CommandFlag `json:"flags"`
	Interactive bool           `json:"interactive"`
}

type CmdFlag struct {
	Name    string      `yaml:"Name,omitempty" json:"name,omitempty"`
	Value   interface{} `yaml:"Value,omitempty" json:"value,omitempty"`
	Options []string    `yaml:"Options,omitempty" json:"options,omitempty"`
	Prompt  string      `yaml:"Prompt,omitempty" json:"prompt,omitempty"`
}

type Prompt struct {
	Options []string
	Prompt  string
}

var (
	routes = []Route{
		{
			Route:   "/command",
			Handler: handleRequestCommand,
		},
		{
			Route:   "/prompt",
			Handler: handleRequestPrompt,
		},
	}
)

func serve(cmdDef *CommandDefinition, command *cobra.Command, args []string) {
	fmt.Printf("\n serve: cmdDef:  %+v\n", cmdDef.Use)

	for _, r := range routes {
		http.HandleFunc(r.Route, r.Handler(cmdDef, command, args))
		log.Infof("Registered route: %s\n", r.Route)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", defaultListenPort), nil))
}

func handleRequestCommand(cmdDef *CommandDefinition, command *cobra.Command, args []string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		resp, _ := json.Marshal(CommandableResponse{
			Cmd:         command.Name(),
			Args:        args,
			Flags:       cmdDef.Flags,
			Interactive: cmdDef.Interactive,
		})

		_, err := fmt.Fprintf(w, "%s", string(resp))
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}

func handleRequestPrompt(
	cmdDef *CommandDefinition,
	command *cobra.Command,
	args []string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			log.Fatalf("Error: %+v", err)
		}

		var data Prompt
		err = json.Unmarshal(buf.Bytes(), &data)
		if err != nil {
			log.Fatalf("Error: %+v", err)
		}

		value, err := promptForStringOption(data.Prompt, data.Options)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		_, err = fmt.Fprintf(w, "%s", value)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}
