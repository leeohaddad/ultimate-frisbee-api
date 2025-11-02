package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheckHandlerV1() echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		requestContext := echoContext.Request().Context()
		err := healthCheck(requestContext)

		if err != nil {
			stringResponse := fmt.Sprintf("unhealthy: %s", err.Error())

			return DispatchEchoResponseFromString(echoContext, http.StatusInternalServerError, stringResponse)
		}

		return DispatchEchoResponseFromString(echoContext, http.StatusOK, "healthy")
	}
}

func healthCheck(context context.Context) error {
	// TODO(lhaddad): configure this health check when postgres is integrated
	// _, err := db.ExecContext(context, "SELECT * FROM teams")
	// if err != nil {
	// 	return fmt.Errorf("could not execute query in database: %w", err)
	// }

	return nil
}
