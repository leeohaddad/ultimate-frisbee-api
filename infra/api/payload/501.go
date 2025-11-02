package payload

import (
	"github.com/labstack/echo/v4"
)

// ValidateNotImplementedPayload validates the params of a request not yet implemented.
func ValidateNotImplementedPayload() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
