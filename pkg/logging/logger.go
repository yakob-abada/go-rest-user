package logging

import "context"

// Logger interface for structured logging
type Logger interface {
	Info(ctx context.Context, msg string, fields map[string]interface{})
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	Error(ctx context.Context, msg string, fields map[string]interface{})
}
