package errorhandler

import (
	"context"

	"github.com/labstack/echo/v4"
)

// ErrorHandler defines an interface for handling API errors
type ErrorHandler interface {
	HandleBadRequest(ctx context.Context, c echo.Context, message string, details map[string]interface{}) error
	HandleInternalServerError(ctx context.Context, c echo.Context, message string, details map[string]interface{}) error
}
