package extensions

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	defaultListenPort = 8080
)

type CommandableResponse struct {
	Cmd   string
	Args  []string
	Flags map[string]string
}

func onCommandRequest(command *cobra.Command, args []string) string {

	flags := map[string]string{}
	command.Flags().VisitAll(func(f *pflag.Flag) {
		// Need more checks here, but making moves
		switch f.Value.Type() {
		case "string":
			flags[f.Name] = f.Value.String()
		}
	})

	returnObj := CommandableResponse{
		Cmd:   command.Name(),
		Args:  args,
		Flags: flags,
	}

	resp, _ := json.Marshal(returnObj)

	return string(resp)
}

func serve(command *cobra.Command, args []string) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		result := onCommandRequest(command, args)

		fmt.Fprintf(w, "%s", result)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", defaultListenPort), nil))
}
