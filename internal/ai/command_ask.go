package ai

import (
	"bytes"
	"encoding/json"
	"github.com/briandowns/spinner"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"time"
)

const ASK_NR_AI_URL = "https://ask-nr-ai.staging-service.nr-ops.net/chat-completion"

var (
	question string
)

type Data struct {
	Assistant string    `json:"assistant"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content""`
}

var cmdAsk = &cobra.Command{
	Use:   "ask",
	Short: "Ask New Relic AI a question",
	Long: `Ask New Relic AI a question

`,
	Example: `newrelic ai ask -q "What is the last version of the NewRelic python agent?"`,
	Run: func(cmd *cobra.Command, args []string) {
		client := &http.Client{}

		accountID := configAPI.RequireActiveProfileAccountID()

		apiKey := configAPI.GetActiveProfileString(config.APIKey)

		if apiKey == "" {
			log.Fatal("an API key is required, set one in your default profile or use the NEW_RELIC_API_KEY environment variable")
		}

		reqBody, _ := json.Marshal(Data{
			Assistant: "grok",
			Messages: []Message{{
				Role:    "user",
				Content: question,
			},
			},
		})

		req, err := http.NewRequest(http.MethodPost, ASK_NR_AI_URL, bytes.NewBuffer(reqBody))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("api-key", apiKey)
		req.Header.Add("x-account-id", string(accountID))

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		s := spinner.New(spinner.CharSets[39], 100*time.Millisecond)
		s.Prefix = "🤖 Loading... 🤖  "
		s.Suffix = "  🤖 Loading... 🤖"
		s.Start()
		time.Sleep(4 * time.Second)
		s.Stop()

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		data := &Data{}
		derr := json.Unmarshal(body, &data)
		if derr != nil {
			log.Fatal(derr)
		}

		for _, m := range data.Messages {
			pterm.DefaultCenter.Println(m.Content)
		}
	},
}

func init() {
	Command.AddCommand(cmdAsk)
	cmdAsk.Flags().StringVarP(&question, "question", "q", "", "a question for New Relic AI in string format")
	utils.LogIfError(cmdAsk.MarkFlagRequired("question"))
}
