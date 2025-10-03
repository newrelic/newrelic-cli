package migrate

import (
	"fmt"
	"strings"
	"time"
)

// printAccountMismatchWarningCI displays a bold warning about account ID mismatches
func printAccountMismatchWarningCI(envAccountID string, mismatchedResources []string) {
	fmt.Printf("\n%s%s%s%s\n", ColorYellow, ColorBold, AccountMismatchWarning, ColorReset)
	fmt.Printf("%sEnvironment NEW_RELIC_ACCOUNT_ID: %s%s\n", ColorBold, envAccountID, ColorReset)
	fmt.Printf("%sThe following resources have different account IDs:%s\n", ColorBold, ColorReset)
	for _, resource := range mismatchedResources {
		fmt.Printf("  - %s%s%s\n", ColorBold, resource, ColorReset)
	}
	fmt.Printf("%sThis may cause import failures. Please verify your account configuration.%s\n\n", ColorYellow, ColorReset)
}

// showStepHeaderCI displays a colored step header
func showStepHeaderCI(step string) {
	fmt.Printf("%s%s%s\n", ColorCyan, step, ColorReset)
}

// showLoadingAnimationCI displays an animated loading message
func showLoadingAnimationCI(message string, duration time.Duration) {
	fmt.Printf("%s", message)

	interval := duration / time.Duration(len(LoadingDots)*3) // Show each dot pattern 3 times

	for i := 0; i < len(LoadingDots)*3; i++ {
		fmt.Printf("\r%s%s   ", message, LoadingDots[i%len(LoadingDots)])
		time.Sleep(interval)
	}

	fmt.Printf("\r%s... ✅\n\n", message)
}

// printCIRecommendationsCI displays the comprehensive CI/CD migration recommendations
func printCIRecommendationsCI(config *ToolConfig, dropRuleData *DropRuleInput, workspacePath string) {
	// Print attention-grabbing separator
	separator := strings.Repeat(SeparatorChar, SeparatorLength)
	fmt.Printf("\n%s\n", separator)
	fmt.Println(MigrationRecommendationsTitle)
	fmt.Printf("%s\n\n", separator)

	// Show loading animation while "preparing recommendations"
	showLoadingAnimationCI("Preparing migration recommendations", 2*time.Second)

	fmt.Println(SuccessMessage)
	fmt.Printf("\nNext steps for your CI/CD pipeline migration (using %s):\n", config.DisplayName)
	fmt.Println()

	time.Sleep(1 * time.Second)

	// Step 1: Copy files to CI
	showStepHeaderCI(fmt.Sprintf("1. COPY GENERATED FILES TO YOUR CI ENVIRONMENT (%s):", config.DisplayName))
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("   📁 Copy the following files to your CI %s workspace:\n", config.DisplayName)
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("      - %s        (Pipeline Cloud Rules configuration)\n", PipelineCloudRulesFile)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("      - %s     (Import blocks)\n", ImportConfigFile)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("      - %s    (Removed blocks for drop rules)\n", RemovedBlocksFile)
	fmt.Println()
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   ⚠️  REQUIREMENT: Ensure %s version >= %s in your CI environment\n", config.DisplayName, MinRequiredVersion)
	fmt.Println("       for import block support.")
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)

	// Step 2: Comment out drop rules and add removals
	showStepHeaderCI("2. PREPARE YOUR EXISTING CI CONFIGURATION:")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("   📝 Comment out ALL existing NRQL drop rule resources in your CI configuration")
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   📝 Copy all content from %s into your CI %s configuration\n", RemovedBlocksFile, config.DisplayName)
	fmt.Println()
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   ⚠️  REQUIREMENT: Ensure %s version >= %s in your CI environment\n", config.DisplayName, MinRemovedBlocksVersion)
	fmt.Println("       for removed block support.")
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)

	// Step 3: Alternative approach
	showStepHeaderCI("3. ALTERNATIVE: MANUAL STATE REMOVAL (if you prefer not to use removed blocks):")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("   If you prefer to manually remove drop rules from state without using")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("   removed blocks, run this command in your CI environment:")
	fmt.Println()

	time.Sleep(500 * time.Millisecond)

	// Generate state rm command
	var resourceNames []string
	for _, resource := range dropRuleData.DropRuleResourceIDs {
		resourceNames = append(resourceNames, fmt.Sprintf("newrelic_nrql_drop_rule.%s", resource.Name))
	}

	stateRmCommand := fmt.Sprintf("   %s state rm %s", config.ToolName, strings.Join(resourceNames, " "))
	fmt.Println(stateRmCommand)
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)

	// Step 4: Final steps
	showStepHeaderCI("4. EXECUTE IN YOUR CI ENVIRONMENT:")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("   After copying files and preparing configuration:")
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   📋 %s plan    (review the migration plan)\n", config.ToolName)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("   📋 %s apply   (execute the migration)\n", config.ToolName)
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)

	// Step 5: Verification
	showStepHeaderCI("5. VERIFICATION:")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("   ✅ Verify that Pipeline Cloud Rules are created successfully")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   ✅ Verify that old NRQL drop rules are removed from state (not destroyed)")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   ✅ Test that your log filtering continues to work as expected")
	fmt.Println()

	time.Sleep(1 * time.Second)

	fmt.Printf("%s\n", separator)
	showLoadingAnimationCI("Finalizing setup", 1*time.Second)
	fmt.Println(MigrationCompleteMessage)
	fmt.Printf("Generated files are ready in: %s\n", workspacePath)
	fmt.Printf("%s\n\n", separator)
}
