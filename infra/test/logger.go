package test

import (
	"testing"

	"github.com/leeohaddad/ultimate-frisbee-api/domain/port/logger"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/adapter/logger/zap"
	"github.com/stretchr/testify/require"
)

var internalLogger logger.Logger

func GetLogger(t *testing.T) logger.Logger {
	t.Helper()

	if internalLogger != nil {
		return internalLogger
	}

	internalLogger, err := zap.NewLogger()
	require.NoError(t, err, "Error when creating test logger")

	return internalLogger
}
