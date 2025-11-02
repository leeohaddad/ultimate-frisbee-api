package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// NotImplementedHandler is the entry point to the application's logic of not implemented stuff.
func NotImplementedHandler() echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		err := echoContext.String(http.StatusOK, "Not Implemented")

		if err != nil {
			return fmt.Errorf("failed to send http not implemented status: %w", err)
		}

		return nil
	}
}
