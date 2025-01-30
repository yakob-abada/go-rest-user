package common

import (
	"context"
	"os"

	"github.com/google/uuid"
)

const CorrelationID = "correlation_id"

func GetEnv(key string, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func GetCorrelationID(ctx context.Context) string {
	correlationID := ""
	if ctx == nil {
		return correlationID
	}

	cID := ctx.Value(CorrelationID)
	switch c := cID.(type) {
	case string:
		correlationID = c
	case uuid.UUID:
		correlationID = c.String()
	}
	return correlationID
}
