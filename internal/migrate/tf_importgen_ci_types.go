package migrate

// DropRuleResource represents a single drop rule resource with its metadata
type DropRuleResource struct {
	Name                      string `json:"name"`
	ID                        string `json:"id"`
	PipelineCloudRuleEntityID string `json:"pipeline_cloud_rule_entity_id"`
}

// DropRuleInput represents the input structure containing drop rule resources
type DropRuleInput struct {
	DropRuleResourceIDs []DropRuleResource `json:"drop_rule_resource_ids"`
}

// ToolConfig holds configuration for the Terraform/OpenTofu tool being used
type ToolConfig struct {
	UseTofu     bool
	ToolName    string
	DisplayName string
}
