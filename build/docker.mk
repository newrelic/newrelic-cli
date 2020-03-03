#
# Makefile fragment for Docker actions
#
DOCKER_ORG        ?= newrelic
DOCKER_IMAGE_NAME ?= cli
DOCKER            ?= docker

docker-run: docker-image
	@echo "=== $(PROJECT_NAME) === [ docker-run       ]: running container:"
	$(DOCKER) run -it $(DOCKER_ORG)/$(DOCKER_IMAGE_NAME):$(PROJECT_VER)


#
# Image management
#
docker-image: compile-linux
	@echo "=== $(PROJECT_NAME) === [ docker-image     ]: building docker image:"
	$(DOCKER) build -t $(DOCKER_ORG)/$(DOCKER_IMAGE_NAME):$(PROJECT_VER) .


docker-push: docker-image docker-login
	@echo "=== $(PROJECT_NAME) === [ docker-push      ]: pushing container to Docker Hub..."
ifneq ($(PROJECT_VER), $(PROJECT_VER_TAGGED))
	@echo
	@echo "Not pushing a dirty version, please make sure you are on the latest release"
	@echo
else
	$(DOCKER) push $(DOCKER_ORG)/$(DOCKER_IMAGE_NAME):$(PROJECT_VER)
endif


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


#
# Cleanup
#
docker-clean: docker-rm docker-rmi

docker-rm:
	@echo "=== $(PROJECT_NAME) === [ docker-clean     ]: removing docker containers:"
	@for i in "$$($(DOCKER) ps -a | grep "$(DOCKER_ORG)/$(DOCKER_IMAGE_NAME)" | cut -d' ' -f 1)"; do \
		if [ ! -z "$$i" ]; then \
			echo "=== $(PROJECT_NAME) === [ docker-clean     ]:     "; \
			$(DOCKER) rm -f $$i ; \
		fi ; \
	done

# Must cleanup containers first
docker-rmi: docker-rm
	@echo "=== $(PROJECT_NAME) === [ docker-clean     ]: removing docker images:"
	@for i in "$$($(DOCKER) images | grep "$(DOCKER_ORG)/$(DOCKER_IMAGE_NAME)" | tr -s ' ' | cut -d' ' -f 3)"; do \
		if [ ! -z "$$i" ]; then \
			echo "=== $(PROJECT_NAME) === [ docker-clean     ]:     $$i"; \
			$(DOCKER) rmi -f $$i ; \
		fi ; \
	done

.PHONY: docker-push docker-login docker-image docker-clean docker-rm docker-rmi docker-run
