# Ultimate Frisbee API

API designed to manage information about Ultimate Frisbee teams, players and tournaments.

## Prerequisites

Before running the project, ensure you have the following installed:

* **Go** (version 1.20 or higher) - [Download here](https://go.dev/dl/)
* **Docker** and **Docker Compose** - [Download Docker Desktop](https://www.docker.com/products/docker-desktop/)
* **Make** - See installation instructions below

### Installing Make on Windows

#### Manual Installation (GnuWin32)

1. Download GnuWin32 Make from [GnuWin32 website](http://gnuwin32.sourceforge.net/packages/make.htm)
2. Run the installer (typically installs to `C:\Program Files (x86)\GnuWin32`)
3. **Add Make to your PATH:**
   * Press `Win + R` and type: `sysdm.cpl`
   * Click the **"Advanced"** tab
   * Click **"Environment Variables"**
   * Under **"System variables"**, find and select `Path`
   * Click **"Edit"**
   * Click **"New"**
   * Add the path: `C:\Program Files (x86)\GnuWin32\bin` (or wherever you installed it)
   * Click **"OK"** on all windows
4. **Restart PowerShell** (important!)
5. Verify installation:

```powershell
make --version
```

## How to run locally

After cloning the repository:

* Run `make setup` to download the project dependencies
* Run `make deps/start` to run the dependencies (postgres, redis) in your local machine
* Run `make db/migration/up` to run the database migrations locally
* Run `make db/seed` to populate the local database with seed data located at `/infra/database/seeds`
* Run `make run/api` to run the Ultimate Frisbee API. You can also run `run/api/watch` to run the api on watch mode (rebuilding and restarting whenever a change is made to a .go file)

**Note for Windows users:** Once Make is installed (see Prerequisites above), all commands work the same way on Windows, Linux, and macOS.

### Testing

* Run `make test` to run all tests on the code base
* Run `make test/unit` to run unit tests
* Run `make test/integration` to run integration tests
* Run `make test/nocache` to clean test cache and run all tests
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

Then open your browser at `http://localhost:9999/#/` to view the OpenAPI documentation (A.K.A. Swagger).

## Postman Collection

A Postman collection is available for testing the API endpoints. Import the file located at:

```text
infra/api/Ultimate Frisbee API.postman_collection.json
```

The collection includes pre-configured requests for:

* Health checks
* Team CRUD operations
* People management

The base URL is already set to `http://127.0.0.1:42007`
