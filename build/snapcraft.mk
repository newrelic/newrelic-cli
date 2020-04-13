#
# Makefile fragment for Snapcraft actions
#
SNAPCRAFT         ?= snapcraft

snapcraft-login:
	@echo "=== $(PROJECT_NAME) === [ snapcraft-login     ]: logging into snapcraft hub"
	@if [ -z "${SNAPCRAFT_TOKEN}" ]; then \
		echo "Failure: SNAPCRAFT_TOKEN not set" ; \
		exit 1 ; \
	fi
	@echo "=== $(PROJECT_NAME) === [ snapcraft-login     ]: using env SNAPCRAFT_TOKEN"
	@echo ${SNAPCRAFT_TOKEN} | $(SNAPCRAFT) login --with -


.PHONY: snapcraft-login
