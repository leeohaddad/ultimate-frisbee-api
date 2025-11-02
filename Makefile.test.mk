################################################################################
## Test related targets
################################################################################

test/nocache: ## Run all tests ignoring cache
	@go clean -testcache
	@make test

test: test/unit test/integration ## Run all tests on our codebase

test/unit: ## Run unit tests
	@go test ./... --tags=unit --parallel 10

test/integration: ## Run integration tests
	@go test ./... --tags=integration --parallel 10
