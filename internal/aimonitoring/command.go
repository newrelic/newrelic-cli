package aimonitoring

import (
	"github.com/spf13/cobra"
)

// Command represents the aimonitoring command
var Command = &cobra.Command{
	Use:   "aimonitoring",
	Short: "Interact with New Relic AI Monitoring",
	Long: `Interact with New Relic AI Monitoring

AI Monitoring telemetry is reported by APM agents instrumenting LLM
libraries (for example OpenAI, Bedrock, LangChain) as NRDB events:
LlmChatCompletionSummary, LlmChatCompletionMessage, LlmEmbedding,
LlmFeedbackMessage, LlmTool, LlmAgent, LlmVectorSearch, and
LlmVectorSearchResult.
`,
	Example: "newrelic aimonitoring application search",
}
