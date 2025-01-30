package logging

import (
	"context"
	"github.com/yakob-abada/go-rest-user/pkg/common"
	"os"

	"github.com/rs/zerolog"
)

// ZeroLogger implements the Logger interface using zerolog
type ZeroLogger struct {
	logger zerolog.Logger
}

// NewZeroLogger initializes a new instance of ZeroLogger
func NewZeroLogger() *ZeroLogger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &ZeroLogger{logger: logger}
}

// Info logs an informational message with structured fields
func (zl *ZeroLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	event := zl.logger.Info()
	event.Str("correlation_id", common.GetCorrelationID(ctx)) // Add correlation ID
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

// Warn logs a warning message with structured fields
func (zl *ZeroLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	event := zl.logger.Warn()
	event.Str("correlation_id", common.GetCorrelationID(ctx))
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

// Error logs an error message with structured fields
func (zl *ZeroLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	event := zl.logger.Error()
	event.Str("correlation_id", common.GetCorrelationID(ctx))
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}
