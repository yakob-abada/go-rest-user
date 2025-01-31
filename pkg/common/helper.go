package common

import (
	"context"
	"github.com/google/uuid"
)

const CorrelationID = "correlation_id"

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
