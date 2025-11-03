# Ultimate Frisbee API

API designed to manage information about Ultimate Frisbee teams, players and tournaments.

## How to run locally

After cloning the repository:
* Run `make setup` to download the project dependencies
* Run `make deps/start` to run the dependencies (postgres, redis) in your local machine
* Run `make db/migration/up` to run the database migrations locally
* Run `make db/seed` to populate the local database with seed data located at `/infra/database/seeds`
* Run `make run/api` to run the Ultimate Frisbee API. You can also run `run/api/watch` to run the api on watch mode (rebuilding and restarting whenever a change is made to a .go file)
 
### Testing 

* Run `make test` to run all tests on the code base
* Run `make test/unit` to run unit tests
* Run `make test/integration` to run integration tests

* Run `make mocks` if you need to update the mock files for an interface

### End-to-End (E2E) Tests

End-to-end tests are located in the `e2e/` package and cover the full API workflow, including health checks, team CRUD operations, error handling, and data validation.

To run the e2e tests:

* Run `go test ./e2e` to execute all e2e tests.
* Alternatively, use the shell script `e2e/run_e2e_tests.sh` for custom or automated e2e test runs.

The e2e tests require the API and its dependencies (Postgres, Redis) to be running. Make sure to follow the setup steps above before running e2e tests.

## How to open API docs

To see the API specification of the Ultimate Frisbee API run the following command:
```sh
make open-api/docs
```

As a result a browser window should be open on `http://localhost:9999/#/` with the OpenAPI docs (A.K.A. Swagger).
