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

// printDuplicateNamesWarning displays a warning about duplicate resource names
func printDuplicateNamesWarning(duplicates map[string][]int) {
	fmt.Printf("\n%s%s⚠️  WARNING: DUPLICATE RESOURCE NAMES DETECTED%s\n", ColorYellow, ColorBold, ColorReset)
	fmt.Printf("%sThe following resource names appear multiple times in the input data:%s\n", ColorBold, ColorReset)

	for name, indices := range duplicates {
		fmt.Printf("  - %s%s%s (appears %d times at positions: %v)\n",
			ColorBold, name, ColorReset, len(indices), indices)
	}

	fmt.Printf("\n%sRESOLUTION:%s Duplicate names will be automatically renamed with random suffixes\n", ColorBold, ColorReset)
	fmt.Printf("to ensure unique resource definitions in the generated Pipeline Cloud Rules configuration.\n")
	fmt.Printf("%sThis prevents Terraform resource conflicts during import.%s\n\n", ColorYellow, ColorReset)
}

// printExecutionHeaderCI displays the execution start header
func printExecutionHeaderCI() {
	separator := strings.Repeat(SeparatorChar, SeparatorLength)
	fmt.Printf("\n%s\n", separator)
	fmt.Println("NEW RELIC CLI: PIPELINE CLOUD RULES MIGRATION (CI/CD)")
	fmt.Printf("%s\n", separator)
	fmt.Println("🚀 Starting Terraform/OpenTofu import configuration generation...")
	fmt.Println("📋 This command will prepare your workspace for CI/CD migration")
	fmt.Printf("%s\n\n", separator)

	// aesthetic sleep interval added
	time.Sleep(time.Millisecond * 1500)
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
func printCIRecommendationsCI(config *ToolConfig, dropRuleData *DropRuleInput, workspacePath string, includeRemovedBlocks bool) {
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
	if includeRemovedBlocks {
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("      - %s    (Removed blocks for drop rules)\n", RemovedBlocksFile)
	}
	fmt.Println()
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   ⚠️  REQUIREMENT: Ensure %s version >= %s in your CI environment\n", config.DisplayName, MinRequiredVersion)
	fmt.Println("       for import block support.")
	if includeRemovedBlocks {
		fmt.Printf("   ⚠️  REQUIREMENT: Ensure %s version >= %s in your CI environment\n", config.DisplayName, MinRemovedBlocksVersion)
		fmt.Println("       for removed block support. To skip this constraint, you may alternatively")
		fmt.Println("       use the `terraform state rm` command to delist drop rule resources.")
	}
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)

	// Step 2: Comment out drop rules (only if using removed blocks)
	stepNumber := 2
	if includeRemovedBlocks {
		showStepHeaderCI("2. PREPARE YOUR EXISTING CI CONFIGURATION:")
		time.Sleep(500 * time.Millisecond)
		fmt.Println("   📝 Comment out ALL existing NRQL drop rule resources in your CI configuration")
		time.Sleep(300 * time.Millisecond)
		fmt.Println("   💡 The removed blocks added in removals.tf will handle the state cleanup automatically")
		fmt.Println("      when you run terraform apply. Alternatively, you may use the `terraform state rm`")
		fmt.Println("      command to delist drop rule resources.")
		fmt.Println()

		time.Sleep(1500 * time.Millisecond)
		stepNumber = 3
	}

	// Step 3/2: Execute terraform plan and apply
	showStepHeaderCI(fmt.Sprintf("%d. EXECUTE IN YOUR CI ENVIRONMENT:", stepNumber))
	time.Sleep(500 * time.Millisecond)
	fmt.Println("   After copying files and preparing configuration:")
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   📋 %s plan    (review the migration plan)\n", config.ToolName)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("   📋 %s apply   (execute the migration)\n", config.ToolName)
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)
	stepNumber++

	// Step 4/3: Manual state removal (only if NOT using removed blocks)
	if !includeRemovedBlocks {
		showStepHeaderCI(fmt.Sprintf("%d. REMOVE EXISTING DROP RULES FROM STATE:", stepNumber))
		time.Sleep(500 * time.Millisecond)
		fmt.Println("   After the import is successful, remove the old NRQL drop rules from state")
		time.Sleep(200 * time.Millisecond)
		fmt.Println("   to complete the migration:")
		fmt.Println()

		time.Sleep(500 * time.Millisecond)

		fmt.Println("   💡 Enabling your CI/CD environment to allow listing resources in")
		fmt.Println("      the Terraform state (using `terraform state list`) can help")
		fmt.Println("      identify drop rule resources to be deleted using the `terraform state rm` command.")
		fmt.Printf("      %s state list | grep nrql_drop_rule\n", config.ToolName)
		fmt.Println()
		time.Sleep(300 * time.Millisecond)

		fmt.Println("   💡 The IDs of drop rule resources specified in the JSON may also be")
		fmt.Println("      leveraged to perform a `state rm`. E.g. `terraform state rm $(terraform state list -id=<drop-rule-ID>)`")

		fmt.Println()

		time.Sleep(1500 * time.Millisecond)
		stepNumber++
	}

	// Final step: Verification
	showStepHeaderCI(fmt.Sprintf("%d. VERIFICATION:", stepNumber))
	time.Sleep(500 * time.Millisecond)
	fmt.Println("   ✅ Verify that Pipeline Cloud Rules are created successfully")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   ✅ Verify that old NRQL drop rules are removed from state (not destroyed)")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   ✅ Test that your log filtering continues to work as expected, via Pipeline Cloud Rules")
	fmt.Println()

	time.Sleep(1 * time.Second)

	fmt.Printf("%s\n", separator)
	showLoadingAnimationCI("Finalizing setup", 1*time.Second)
	fmt.Println(MigrationCompleteMessage)
	fmt.Printf("Generated files are ready in: %s\n", workspacePath)
	fmt.Printf("%s\n\n", separator)
}
