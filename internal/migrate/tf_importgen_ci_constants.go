package migrate

// Tool names and display names
const (
	TerraformToolName = "terraform"
	TofuToolName      = "tofu"
	TerraformDisplay  = "Terraform"
	OpenTofuDisplay   = "OpenTofu"
)

// Version requirements
const (
	MinRequiredVersion      = "1.5"
	MinRemovedBlocksVersion = "1.7"
)

// File names
const (
	ProviderConfigFile     = "provider.tf"
	ImportConfigFile       = "imports.tf"
	RemovedBlocksFile      = "removals.tf"
	PipelineCloudRulesFile = "pcrs.tf"
)

// Environment variables
var RequiredEnvVars = []string{"NEW_RELIC_ACCOUNT_ID", "NEW_RELIC_API_KEY"}

const OptionalRegionEnvVar = "NEW_RELIC_REGION"

// File patterns for workspace validation
var ConflictingFilePatterns = []string{
	".tfstate",
	".tfstate.backup",
	".terraform",
	ProviderConfigFile,
	ImportConfigFile,
	RemovedBlocksFile,
	PipelineCloudRulesFile,
}

// Provider configuration template
const ProviderConfigTemplate = `terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
      version = "~> 3.0"
    }
  }
}

provider "newrelic" {
  # Configuration will be taken from environment variables:
  # NEW_RELIC_API_KEY
  # NEW_RELIC_ACCOUNT_ID
  # NEW_RELIC_REGION (optional, defaults to US)
}
`

// Import block template
const ImportBlockTemplate = `import {
  to = newrelic_pipeline_cloud_rule.%s
  id = "%s"
}`

// Removed block template
const RemovedBlockTemplate = `removed {
  from = newrelic_nrql_drop_rule.%s

  lifecycle {
    destroy = false
  }
}`

// UI constants
const (
	SeparatorLength = 80
	SeparatorChar   = "="
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorCyan   = "\033[1;36m"
	ColorYellow = "\033[1;33m"
)

// Loading animation
var LoadingDots = []string{".", "..", "...", "...."}

// Messages
const (
	MigrationRecommendationsTitle = "IMPORTANT: CI/CD MIGRATION RECOMMENDATIONS"
	SuccessMessage                = "✅ Local workspace setup completed successfully!"
	MigrationCompleteMessage      = "Migration workspace prepared successfully! 🎉"
	AccountMismatchWarning        = "⚠️  WARNING: ACCOUNT ID MISMATCH DETECTED"
)
