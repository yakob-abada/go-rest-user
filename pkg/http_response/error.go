package http_response

import (
	"github.com/google/uuid"
	"time"
)

type Error struct {
	CorrelationId uuid.UUID `json:"correlation_id"`
	Message       string    `json:"message"`
	Status        int       `json:"status"`
	Time          time.Time `json:"time"`
}
