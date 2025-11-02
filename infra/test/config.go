package test

import (
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/config"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/adapter/config/viper"
	"github.com/stretchr/testify/require"
)

var internalConfig *config.Application

func GetTestConfiguration(t *testing.T) *config.Application {
	t.Helper()

	if internalConfig != nil {
		return internalConfig
	}

	// This is a pretty nasty workaround to get the root folder =(
	// see this for more details: https://stackoverflow.com/questions/23847003/golang-tests-and-working-directory
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")

	viperConfig, err := viper.NewConfig(fmt.Sprintf("%s/config/test.yaml", rootDir))
	require.NoError(t, err, "Error when creating test config")

	internalConfig, err = viperConfig.GetConfigs()
	require.NoError(t, err, "Error when creating test config")

	return internalConfig
}
