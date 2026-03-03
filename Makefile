# ===================================================================
# Vairables for system
# ===================================================================
SYSTEM_NAME = ABCD
GO_CMD = go
KEY_SERVER_MODE = SYSTEM_MODE
VALUE_SERVER_MODE_DEV = dev
VALUE_SERVER_MODE_PROD = prod
DOCKER_CMD = docker
DOCKER_IMAGE = $(SYSTEM_NAME)
DOCKER_IMAGE_SAFE = $(shell echo $(DOCKER_IMAGE) | tr '[:upper:]' '[:lower:]')
DOCKER_TAG = latest
DOCKER_PLATFORM = linux/amd64

# ===================================================================
# Define commands using
# ===================================================================
# Global
all: help
.PHONY: help
# For dockers
.PHONY: docker_run_dev docker_stop_dev docker_rm_dev docker_run_prod docker_stop_prod docker_rm_prod
# For build
.PHONY: build docker_build
# For development
.PHONY: dev_server

# ===================================================================
# Deploy commands
# ===================================================================

dev_server: ## Run development server
	@echo "Running development server ..."
	$(KEY_SERVER_MODE)=$(VALUE_SERVER_MODE_DEV) $(GO_CMD) run ./cmd/server/main.go

build: ## Build server binary
	@echo "Building server binary ..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build -trimpath -ldflags="-s -w -buildid=" -o ./bin/server ./cmd/server/main.go

docker_build: ## Build docker image
	@echo "Building docker image $(DOCKER_IMAGE_SAFE):$(DOCKER_TAG) ..."
	$(DOCKER_CMD) buildx build --platform $(DOCKER_PLATFORM) --load -f Dockerfile -t $(DOCKER_IMAGE_SAFE):$(DOCKER_TAG) .

docker_rm_prod: ## Remove docker compose for production environment
	@echo "Todo remove docker compose for production environment"

docker_stop_prod: ## Stop docker compose for production environment
	@echo "Stopping docker compose for production environment ..."
	docker compose -f environment/docker-compose.yml -p $(SYSTEM_NAME)-prod stop

docker_run_prod: ## Run docker compose for production environment
	@echo "Running docker compose for production environment ..."
	docker compose -f environment/docker-compose.yml -p $(SYSTEM_NAME)-prod up -d

docker_rm_dev: ## Remove docker compose for development environment
	@echo "Todo remove docker compose for development environment"

docker_stop_dev: ## Stop docker compose for development environment
	@echo "Stopping docker compose for development environment ..."
	docker compose -f environment/docker-compose-dev.yml -p $(SYSTEM_NAME)-dev stop

docker_run_dev: ## Run docker compose for development environment
	@echo "Running docker compose for development environment ..."
	docker compose -f environment/docker-compose-dev.yml -p $(SYSTEM_NAME)-dev up -d

help: ## Help commands used in this Makefile
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'