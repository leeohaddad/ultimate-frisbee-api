################################################################################
## Database related targets
################################################################################

PG_HOST := localhost:42013
DATABASE_CONNECTION_STRING := postgres://ultimate_frisbee_manager_user:some_password@$(PG_HOST)/ultimate_frisbee_manager?sslmode=disable

run/migration: lint ## Run the migration tool
	@go run *.go -e migration

services/database/up: ## Starts a database container for local environment
	@docker compose up -d postgres

services/database/down: ## Stops the database container for local environment
	@docker compose stop postgres

services/database/logs: ## Show the database container logs on the local environment
	@docker compose logs -f postgres

db/migration/create: ensure-migrate-installed ## Given a MIGRATION_NAME, creates a new migration
	@migrate create -seq -ext sql -dir ./infra/database/migrations $(MIGRATION_NAME)

db/migration/up: ensure-migrate-installed ## Run migrations against the local database
	@migrate -path ./infra/database/migrations -database $(DATABASE_CONNECTION_STRING) up

db/migration/down: ensure-migrate-installed ## Rollback the migrations against the local database
	@migrate -path ./infra/database/migrations -database $(DATABASE_CONNECTION_STRING) down

db/drop: ensure-migrate-installed ## Drop local database elements
	@migrate -path ./infra/database/migrations -database $(DATABASE_CONNECTION_STRING) drop

db/seed: ## Populate the local database with seed data using the application
	@go run main.go api.go migration.go -e seed

db/seed/sql: ensure-psql-installed ## Run seed scripts against the local database (alternative method)
	@psql $(DATABASE_CONNECTION_STRING) -f ./infra/database/seeds/00001_teams.sql > /dev/null
	@psql $(DATABASE_CONNECTION_STRING) -f ./infra/database/seeds/00002_people.sql > /dev/null
	@echo Database seeded

ensure-migrate-installed:
	@command -v migrate >/dev/null 2>&1 || { echo >&2 "migrate is necessary to run this command. please run 'brew install golang-migrate' and try again"; exit 1; }

ensure-psql-installed:
	@command -v psql >/dev/null 2>&1 || { echo >&2 "psql is necessary to run this command. please run 'brew install postgres' and try again"; exit 1; }
