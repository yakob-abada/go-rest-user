package bootstrap

import (
	"github.com/yakob-abada/go-rest-user/pkg/handler"
	"gorm.io/gorm"
)

func NewHealthHandler(db *gorm.DB) *handler.HealthHandler {
	return handler.NewHealthHandler(db)
}
