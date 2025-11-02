package handler

import (
	"fmt"

	"github.com/leeohaddad/ultimate-frisbee-api/infra/api/handler/result"

	"github.com/labstack/echo/v4"
)

// DispatchEchoResponseFromHandlerResult dispatches an Echo Response using information fetched from the handler HTTP.
func DispatchEchoResponseFromHandlerResult(echoContext echo.Context, param result.HTTP) error {
	switch param.ResponseType {
	case result.ResponseBodyTypes.JSON:
		err := echoContext.JSON(param.StatusCode, param.JSONResponse)
		if err != nil {
			return fmt.Errorf("failed to write JSON response: %w", err)
		}

		return nil
	case result.ResponseBodyTypes.String:
		return DispatchEchoResponseFromString(echoContext, param.StatusCode, param.StringResponse)
	}

	return fmt.Errorf("unknown response body type: %s", param.ResponseType)
}

func DispatchEchoResponseFromString(echoContext echo.Context, statusCode int, stringResponse string) error {
	err := echoContext.String(statusCode, stringResponse)
	if err != nil {
		return fmt.Errorf("failed to write string response: %w", err)
	}

	return nil
}
