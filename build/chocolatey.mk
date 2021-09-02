#
# Makefile fragment for Chocolatey actions
#
CHOCOLATEY_BUILD_DIR    ?= build/package/chocolatey
CHOCO                   ?= docker run --rm -v $$PWD/$(CHOCOLATEY_BUILD_DIR):$$PWD -w $$PWD linuturk/mono-choco

chocolatey-publish: chocolatey-build
	@echo "=== $(PROJECT_NAME) === [ chocolatey-publish ]: publishing chocolatey package"
	@if [ -z "${CHOCOLATEY_API_KEY}" ]; then \
		echo "Failure: CHOCOLATEY_API_KEY not set" ; \
		exit 1 ; \
	fi
	$(CHOCO) push --source https://chocolatey.org/ -k ${CHOCOLATEY_API_KEY} newrelic-cli.${PROJECT_VER_TAGGED}.nupkg


chocolatey-build:
	@echo "=== $(PROJECT_NAME) === [ chocolatey-build ]: publishing chocolatey package"
	@cp LICENSE $(CHOCOLATEY_BUILD_DIR)/tools/LICENSE.txt
	@curl -sL -o $(CHOCOLATEY_BUILD_DIR)/tools/NewRelicCLIInstaller.msi https://download.newrelic.com/install/newrelic-cli/v${PROJECT_VER_TAGGED}/NewRelicCLIInstaller.msi
	@rm -f newrelic-cli.${PROJECT_VER_TAGGED}.nupkg
	@sed -i.bak -e "s/    <version>.*<\/version>/    <version>${PROJECT_VER_TAGGED}<\/version>/g" $(CHOCOLATEY_BUILD_DIR)/newrelic-cli.nuspec
	$(CHOCO) pack

.PHONY: chocolatey-build chocolatey-publish
