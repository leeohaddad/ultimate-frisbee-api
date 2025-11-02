package test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	databasePostgres "github.com/leeohaddad/ultimate-frisbee-api/infra/database/postgres"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/require"
)

const (
	postgresUser     = "ultimate_frisbee_manager_user"
	postgresPassword = "some_password"
	postgresDatabase = "ultimate_frisbee_manager"
	maxConnections   = 95 // Today the default is 100 on Postgres, so we used a slightly lower value because the default test make the suite instable
)

var (
	internalRootDatabaseClient databasePostgres.Client
	databaseContainer          *gnomock.Container
	dbNumber                   int32
	migrationPath              string
	lock                       sync.Mutex
)

func GetDatabaseClient(t *testing.T) databasePostgres.Client {
	t.Helper()

	lock.Lock()
	defer lock.Unlock()

	rootClient := getRootDatabaseClient(t)

	// TODO(ddias): how we can solve this without explicitly cleaning idle connections?
	// checks if we have too many idle connections before creating a new one
	checkForTooManyIdleConnections(t)

	number := atomic.AddInt32(&dbNumber, 1)
	databaseName := fmt.Sprintf("%s_%d", postgresDatabase, number)

	_, err := rootClient.ExecuteCommand(context.Background(), fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", databaseName, postgresDatabase))
	require.NoError(t, err, "Error when creating test database")

	logger := GetLogger(t)
	client := databasePostgres.NewClient(logger)

	connectionString := buildConnectionString(databaseName)

	err = client.Connect(connectionString)
	require.NoError(t, err, "Error when connecting to test database: %w", err)

	err = databasePostgres.RunMigrations(migrationPath, connectionString)
	require.NoError(t, err, "Error when populating test database: %w", err)

	return client
}

func getRootDatabaseClient(t *testing.T) databasePostgres.Client {
	t.Helper()

	if internalRootDatabaseClient != nil {
		return internalRootDatabaseClient
	}

	logger := GetLogger(t)
	internalRootDatabaseClient = databasePostgres.NewClient(logger)

	connectionString := buildConnectionString(postgresDatabase)

	err := internalRootDatabaseClient.Connect(connectionString)
	require.NoError(t, err, "Error when creating root test database")

	return internalRootDatabaseClient
}

func buildConnectionString(databaseName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		postgresUser, postgresPassword, databaseContainer.Host, databaseContainer.DefaultPort(), databaseName)
}

func checkForTooManyIdleConnections(t *testing.T) {
	t.Helper()

	rootClient := getRootDatabaseClient(t)

	var countData struct{ Counter int32 }
	_, err := rootClient.ExecuteQuery(context.Background(), &countData, "SELECT count(*) as counter FROM pg_stat_activity WHERE state in ('active', 'idle')")
	require.NoError(t, err, "Error when checking connections on test database")

	if countData.Counter <= maxConnections {
		return
	}

	// TODO(ddias): the correct thing to do here is to add a synchronization barrier
	time.Sleep(10 * time.Second) // wait for other tests to finish

	query := `WITH inactive_connections AS (
            SELECT
                pid,
                rank() over (partition by client_addr order by backend_start ASC) as rank
            FROM
                pg_stat_activity
            WHERE
                -- Exclude the thread owned connection (ie no auto-kill)
                pid <> pg_backend_pid( )
            AND
                -- Exclude known applications connections
                application_name !~ '(?:psql)|(?:pgAdmin.+)'
            AND
                -- Include inactive connections only
                state in ('idle', 'idle in transaction', 'idle in transaction (aborted)', 'disabled')
        )
        SELECT
            pg_terminate_backend(pid)
        FROM
            inactive_connections
        WHERE
            rank > 1 -- Leave one connection for each application connected to the database`

	_, err = rootClient.ExecuteCommand(context.Background(), query)
	require.NoError(t, err, "Error when killing idle connections on test database")
}

func StartDatabaseEnvironment(migrationRelativePath string) error {
	p := postgres.Preset(
		postgres.WithUser(postgresUser, postgresPassword),
		postgres.WithDatabase(postgresDatabase),
	)

	container, err := gnomock.Start(p)
	if err != nil {
		return fmt.Errorf("failed to start database environment: %w", err)
	}

	databaseContainer = container
	migrationPath = migrationRelativePath

	return nil
}

func StopDatabaseEnvironment() error {
	err := gnomock.Stop(databaseContainer)

	if err != nil {
		// We don't want to stop the test execution because of a problem with the test database environment
		// So we just log the error and continue
		// return fmt.Errorf("failed to stop test database environment: %w", err)

		fmt.Printf("failed to stop test database environment: %v\n", err)
		fmt.Println("please make sure to clean the containers manually after the test!")
	}

	return nil
}
