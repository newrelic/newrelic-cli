#
# Makefile fragment for Docker actions
#
DOCKER            ?= docker
DOCKER_FILE       ?= build/package/Dockerfile
DOCKER_IMAGE      ?= newrelic/cli
DOCKER_IMAGE_TAG  ?= snapshot

# Build the docker image
docker-build: compile-linux
	@echo "=== $(PROJECT_NAME) === [ docker-build     ]: Creating docker image: $(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG) ..."
	docker build -f $(DOCKER_FILE) -t $(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG) $(BUILD_DIR)/linux/


docker-login:
	@echo "=== $(PROJECT_NAME) === [ docker-login     ]: logging into docker hub"
	@if [ -z "${DOCKER_USERNAME}" ]; then \
		echo "Failure: DOCKER_USERNAME not set" ; \
		exit 1 ; \
	fi
	@if [ -z "${DOCKER_PASSWORD}" ]; then \
		echo "Failure: DOCKER_PASSWORD not set" ; \
		exit 1 ; \
	fi
	@echo "=== $(PROJECT_NAME) === [ docker-login     ]: username: '$$DOCKER_USERNAME'"
	@echo ${DOCKER_PASSWORD} | $(DOCKER) login -u ${DOCKER_USERNAME} --password-stdin


# Push the docker image
docker-push: docker-login docker-build
	@echo "=== $(PROJECT_NAME) === [ docker-push      ]: Pushing docker image: $(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG) ..."
	@$(DOCKER) push $(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG)

.PHONY: docker-build docker-login docker-push
