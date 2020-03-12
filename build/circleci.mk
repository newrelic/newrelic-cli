#
# Makefile fragment for CircleCI Stuff
#
CIRCLE_NAMESPACE ?= newrelic
ORB_NAME         ?= cli

CIRCLE_CMD   ?= circleci
DIST_DIR     ?= ./dist
ORB_DIR      ?= build/package/circleci-orb

circle-clean:
	@echo "=== $(PROJECT_NAME) === [ circle-clean         ]: distribution files..."
	@rm -rfv $(DIST_DIR)/*

circle-orb-pack:
	@echo "=== $(PROJECT_NAME) === [ circle-orb-pack      ]: Creating orb config from '$(ORB_DIR)'"
	@mkdir -p $(DIST_DIR)
	@$(CIRCLE_CMD) config pack $(ORB_DIR) > $(DIST_DIR)/orb.yml

circle-orb-validate: circle-orb-pack
	@echo "=== $(PROJECT_NAME) === [ circle-orb-validate  ]: Validating orb..."
	@$(CIRCLE_CMD) orb validate $(DIST_DIR)/orb.yml

circle-orb-publish: circle-orb-pack circle-orb-validate
	@echo "=== $(PROJECT_NAME) === [ circle-orb-validate  ]: Publishing  orb to dev..."
	@$(CIRCLE_CMD) orb publish $(DIST_DIR)/orb.yml $(CIRCLE_NAMESPACE)/$(ORB_NAME)@dev:$(PROJECT_VER_TAGGED)

circle-orb: circle-clean circle-orb-pack circle-orb-validate circle-orb-publish

.PHONY: circle-clean circle-orb circle-orb-pack circle-orb-publish circle-orb-validate
