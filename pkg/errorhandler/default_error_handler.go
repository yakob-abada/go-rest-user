package errorhandler

import (
	"context"
	"github.com/yakob-abada/go-rest-user/pkg/common"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yakob-abada/go-rest-user/pkg/logging"
)

// DefaultErrorHandler struct manages API error responses
type DefaultErrorHandler struct {
	Logger logging.Logger
}

// NewErrorHandler creates a new instance of DefaultErrorHandler
func NewErrorHandler(logger logging.Logger) *DefaultErrorHandler {
	return &DefaultErrorHandler{Logger: logger}
}

// HandleBadRequest sends a structured bad request response
func (eh *DefaultErrorHandler) HandleBadRequest(ctx context.Context, c echo.Context, message string, details map[string]interface{}) error {
	//eh.logger.Warn(ctx, message, details)
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"error":          message,
		"details":        details,
		"correlation_id": common.GetCorrelationID(ctx),
	})
}

// HandleInternalServerError sends a structured internal server error response
func (eh *DefaultErrorHandler) HandleInternalServerError(ctx context.Context, c echo.Context, message string, details map[string]interface{}) error {
	//eh.logger.Error(ctx, message, details)
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error":          message,
		"details":        details,
		"correlation_id": common.GetCorrelationID(ctx),
	})
}
