include Makefile.*

APP_NAME ?= ultimate-frisbee-manager
DOCKER_COMPOSE_FILE ?= "docker/docker-compose.yaml"

help: Makefile ## Show list of commands
	@echo "Choose a command to run in "$(APP_NAME)":"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

setup/linter: ## Setup linter to be used by CI
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1

setup/dev: setup/linter setup ## Download project dependencies for development mode
	@go get github.com/cespare/reflex
	@go get github.com/golang/mock/mockgen

setup: ## Download project dependencies
	@go mod download
	@go mod tidy

lint: ## Run linter on project
	@golangci-lint run

lint/clean: ## Run linter on project after cleaning terminal (useful when you need to run lint too many times)
	@clear
	@make lint
	@echo "linter finished!"

lint/watch: ## Run linter on project and watch for changes
	@reflex -r '^.*\.go' -s -d none -- make lint/clean

fmt: ## Run go fmt on project files
	@go fmt ./...

build: ## Build the Ultimate Frisbee API
	@go build -o bin/ultimate_frisbee_api *.go

deps/start: ## Start dependencies locally
	docker-compose up -d

deps/stop: ## Stop dependencies locally
	docker-compose down

run/api: lint ## Run the api
	@go run . -e api

run/worker: lint ## Run the worker
	@go run . -e worker

run/api/watch: ## Run the api on watch mode (rebuilding and restarting whenever a change is made to a .go file)
	@reflex -r '^.*\.go' -s -d none -- make run/api

open-api/docs: ## Open OpenAPI documentation for Ultimate Frisbee API
	@docker run -p 9999:8080 -e SWAGGER_JSON=/docs/public-api.json -v "$(PWD)/docs/openapi:/docs" swaggerapi/swagger-ui
	@open http://localhost:9999

test/e2e: ## Run end-to-end tests against the running API
	@./run_e2e_tests.sh