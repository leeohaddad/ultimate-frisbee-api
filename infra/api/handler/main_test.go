package handler_test

import (
	"os"
	"testing"

	"github.com/leeohaddad/ultimate-frisbee-api/infra/test"
)

const migrationPath = "../../database/migrations"

func TestMain(m *testing.M) {
	err := test.StartDatabaseEnvironment(migrationPath)
	if err != nil {
		panic(err)
	}

	exitVal := m.Run()

	err = test.StopDatabaseEnvironment()

	if err != nil {
		panic(err)
	}

	os.Exit(exitVal)
}
