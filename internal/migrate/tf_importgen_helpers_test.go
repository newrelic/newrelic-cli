package migrate

import (
	"strings"
	"testing"
)

func TestValidateResourceIdentifier(t *testing.T) {
	tests := []struct {
		name        string
		identifier  string
		shouldError bool
		description string
	}{
		// Valid resource identifiers
		{
			name:        "valid simple resource",
			identifier:  "newrelic_nrql_drop_rule.my_rule",
			shouldError: false,
			description: "Standard Terraform resource address",
		},
		{
			name:        "valid resource with count index",
			identifier:  "newrelic_nrql_drop_rule.my_rule[0]",
			shouldError: false,
			description: "Resource with count meta-argument",
		},
		{
			name:        "valid resource with for_each key",
			identifier:  "newrelic_nrql_drop_rule.my_rule[\"prod\"]",
			shouldError: false,
			description: "Resource with for_each meta-argument",
		},
		{
			name:        "valid module resource",
			identifier:  "module.drop_rules.newrelic_nrql_drop_rule.rule",
			shouldError: false,
			description: "Resource inside a module",
		},
		{
			name:        "valid nested module resource",
			identifier:  "module.parent.module.child.newrelic_nrql_drop_rule.rule",
			shouldError: false,
			description: "Resource inside nested modules",
		},

		// Invalid/malicious resource identifiers - Command Injection Attempts
		{
			name:        "injection attempt - lock=false flag",
			identifier:  "-lock=false",
			shouldError: true,
			description: "Attempt to inject -lock=false flag to disable state locking",
		},
		{
			name:        "injection attempt - chdir flag",
			identifier:  "-chdir=/tmp",
			shouldError: true,
			description: "Attempt to inject -chdir flag to change working directory",
		},
		{
			name:        "injection attempt - state flag",
			identifier:  "-state=/tmp/evil.tfstate",
			shouldError: true,
			description: "Attempt to inject -state flag to use malicious state file",
		},
		{
			name:        "injection attempt - backup flag",
			identifier:  "-backup=false",
			shouldError: true,
			description: "Attempt to inject -backup flag",
		},
		{
			name:        "injection attempt - lock-timeout flag",
			identifier:  "-lock-timeout=0s",
			shouldError: true,
			description: "Attempt to inject -lock-timeout flag",
		},

		// Edge cases
		{
			name:        "empty identifier",
			identifier:  "",
			shouldError: true,
			description: "Empty string should be rejected",
		},
		{
			name:        "no dot separator",
			identifier:  "invalid_resource",
			shouldError: true,
			description: "Missing dot separator (invalid Terraform address format)",
		},
		{
			name:        "single dash prefix (flag)",
			identifier:  "-resource.name",
			shouldError: true,
			description: "Leading dash should be rejected",
		},
		{
			name:        "double dash prefix (flag)",
			identifier:  "--lock=false",
			shouldError: true,
			description: "Leading double dash should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateResourceIdentifier(tt.identifier)

			if tt.shouldError && err == nil {
				t.Errorf("Expected error for %q (%s), but got none", tt.identifier, tt.description)
			}

			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error for %q (%s), but got: %v", tt.identifier, tt.description, err)
			}
		})
	}
}

func TestGenerateStateRmCommands_WithValidation(t *testing.T) {
	tests := []struct {
		name          string
		toolName      string
		resources     []string
		expectedCount int
		description   string
	}{
		{
			name:     "valid resources",
			toolName: "terraform",
			resources: []string{
				"newrelic_nrql_drop_rule.rule1",
				"newrelic_nrql_drop_rule.rule2",
			},
			expectedCount: 2,
			description:   "All valid resources should generate commands",
		},
		{
			name:     "mixed valid and invalid resources",
			toolName: "terraform",
			resources: []string{
				"newrelic_nrql_drop_rule.rule1",
				"-lock=false",
				"newrelic_nrql_drop_rule.rule2",
			},
			expectedCount: 2,
			description:   "Invalid resources should be skipped",
		},
		{
			name:     "all invalid resources",
			toolName: "terraform",
			resources: []string{
				"-lock=false",
				"-chdir=/tmp",
				"-state=/evil.tfstate",
			},
			expectedCount: 0,
			description:   "All invalid resources should result in no commands",
		},
		{
			name:          "empty resource list",
			toolName:      "terraform",
			resources:     []string{},
			expectedCount: 0,
			description:   "Empty list should result in no commands",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands := generateStateRmCommands(tt.toolName, tt.resources)

			if len(commands) != tt.expectedCount {
				t.Errorf("Expected %d commands (%s), but got %d", tt.expectedCount, tt.description, len(commands))
			}

			// Verify no commands contain injected flags
			for _, cmd := range commands {
				if containsMaliciousFlags(cmd) {
					t.Errorf("Command contains potentially malicious flags: %s", cmd)
				}
			}
		})
	}
}

// Helper function to detect common malicious flags in commands
func containsMaliciousFlags(command string) bool {
	// Check if the command or any part of it starts with a dash (flag)
	if len(command) > 0 && command[0] == '-' {
		return true
	}

	maliciousPatterns := []string{
		" -lock=false",
		" -chdir=",
		" -state=",
		" -backup=false",
		" -lock-timeout=0",
	}

	for _, pattern := range maliciousPatterns {
		if len(command) > len(pattern) {
			for i := 0; i <= len(command)-len(pattern); i++ {
				if command[i:i+len(pattern)] == pattern {
					return true
				}
			}
		}
	}

	return false
}

func TestValidateResourceIdentifier_SecurityScenarios(t *testing.T) {
	// These are the exact attack scenarios from the security report NR-526865
	securityTests := []struct {
		name       string
		identifier string
		attackType string
	}{
		{
			name:       "PoC1 - Disable state locking",
			identifier: "-lock=false",
			attackType: "State locking bypass",
		},
		{
			name:       "PoC2 - Change execution context",
			identifier: "-chdir=/tmp",
			attackType: "Execution context manipulation",
		},
		{
			name:       "PoC3 - State file redirection",
			identifier: "-state=/tmp/evil.tfstate",
			attackType: "State file manipulation",
		},
	}

	for _, tt := range securityTests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateResourceIdentifier(tt.identifier)
			if err == nil {
				t.Errorf("SECURITY FAILURE: %s attack (%q) was not blocked!", tt.attackType, tt.identifier)
				t.Errorf("This attack could affect both tf-delist and tf-update commands")
			} else {
				t.Logf("✓ Successfully blocked %s attack (protects tf-delist and tf-update): %v", tt.attackType, err)
			}
		})
	}
}

func TestGenerateTargetCommands_WithValidation(t *testing.T) {
	tests := []struct {
		name             string
		commandType      CommandType
		resources        []string
		shouldContain    string
		shouldNotContain string
		description      string
	}{
		{
			name:          "tf-update with valid resources",
			commandType:   CommandUpdate,
			resources:     []string{"newrelic_nrql_drop_rule.rule1", "newrelic_nrql_drop_rule.rule2"},
			shouldContain: "-target=",
			description:   "Valid resources should generate target flags",
		},
		{
			name:             "tf-update with injection attempt",
			commandType:      CommandUpdate,
			resources:        []string{"newrelic_nrql_drop_rule.rule1", "-lock=false"},
			shouldContain:    "-target=newrelic_nrql_drop_rule.rule1",
			shouldNotContain: "-lock=false",
			description:      "Injection attempts should be filtered out",
		},
		{
			name:          "tf-delist with valid resources",
			commandType:   CommandDelist,
			resources:     []string{"newrelic_nrql_drop_rule.rule1"},
			shouldContain: "state rm",
			description:   "Valid delist command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &CommandContext{
				ToolConfig: &ToolConfig{
					ToolName:    "terraform",
					DisplayName: "Terraform",
				},
				CommandType: tt.commandType,
			}

			planCmd, actionCmd := ctx.GenerateTargetCommands(tt.resources)

			// Check that commands contain expected content
			if tt.shouldContain != "" {
				combinedCmd := planCmd + " " + actionCmd
				if !strings.Contains(combinedCmd, tt.shouldContain) {
					t.Errorf("Expected commands to contain %q, but didn't find it in:\nPlan: %s\nAction: %s",
						tt.shouldContain, planCmd, actionCmd)
				}
			}

			// Check that malicious content is not present
			if tt.shouldNotContain != "" {
				combinedCmd := planCmd + " " + actionCmd
				if strings.Contains(combinedCmd, tt.shouldNotContain) {
					t.Errorf("Commands should NOT contain %q (%s), but found it in:\nPlan: %s\nAction: %s",
						tt.shouldNotContain, tt.description, planCmd, actionCmd)
				}
			}
		})
	}
}

func TestUpdateCommand_InjectionPrevention(t *testing.T) {
	// Specific tests for tf-update command injection prevention
	attackVectors := []struct {
		name          string
		resourceInput string
		expectedBlock bool
		description   string
	}{
		{
			name:          "Attempt to add malicious -target flag",
			resourceInput: "-target=other.resource",
			expectedBlock: true,
			description:   "Should block attempt to add additional target",
		},
		{
			name:          "Attempt to inject -var flag",
			resourceInput: "-var=account_id=attacker",
			expectedBlock: true,
			description:   "Should block variable injection",
		},
		{
			name:          "Attempt to inject -backend-config",
			resourceInput: "-backend-config=path=/tmp/evil",
			expectedBlock: true,
			description:   "Should block backend config injection",
		},
		{
			name:          "Valid resource should pass",
			resourceInput: "newrelic_nrql_drop_rule.my_rule",
			expectedBlock: false,
			description:   "Legitimate resource should work",
		},
	}

	for _, tt := range attackVectors {
		t.Run(tt.name, func(t *testing.T) {
			err := validateResourceIdentifier(tt.resourceInput)
			blocked := (err != nil)

			if blocked != tt.expectedBlock {
				if tt.expectedBlock {
					t.Errorf("SECURITY FAILURE: %s (%q) was not blocked!", tt.description, tt.resourceInput)
				} else {
					t.Errorf("FALSE POSITIVE: Valid input (%q) was incorrectly blocked: %v", tt.resourceInput, err)
				}
			} else if blocked {
				t.Logf("✓ Successfully blocked: %s", tt.description)
			}
		})
	}
}
