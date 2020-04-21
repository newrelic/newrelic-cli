#
# Makefile fragment for Chocolatey actions
#
CHOCO         ?= docker run --rm -v $$PWD:$$PWD -w $$PWD linuturk/mono-choco
CHOCOLATEY_BUILD_DIR	?= build/package/chocolatey

chocolatey-publish: chocolatey-build
	@echo "=== $(PROJECT_NAME) === [ choco-publish     ]: publishing chocolatey package"
	@if [ -z "${CHOCOLATEY_API_KEY}" ]; then \
		echo "Failure: CHOCOLATEY_API_KEY not set" ; \
		exit 1 ; \
	fi
	@cd $(CHOCOLATEY_BUILD_DIR) && \
		$(CHOCO) push --source https://chocolatey.org/ -k ${CHOCOLATEY_API_KEY} newrelic-cli.${PROJECT_VER_TAGGED}.nupkg \
	; cd -


chocolatey-build:
	@echo "=== $(PROJECT_NAME) === [ choco-build     ]: building chocolatey package"
	@cp LICENSE $(CHOCOLATEY_BUILD_DIR)/tools/LICENSE.txt
	@cd $(CHOCOLATEY_BUILD_DIR) && \
		curl -sL -o NewRelicCLIInstaller.msi https://github.com/newrelic/newrelic-cli/releases/download/v${PROJECT_VER_TAGGED}/NewRelicCLIInstaller.msi && \
		rm -f newrelic-cli.${PROJECT_VER_TAGGED}.nupkg && \
		sed -i '' -e "s/    <version>.*<\/version>/    <version>${PROJECT_VER_TAGGED}<\/version>/g" newrelic-cli.nuspec && \
		$(CHOCO) pack \
	; cd -

.PHONY: chocolatey-build chocolatey-publish
