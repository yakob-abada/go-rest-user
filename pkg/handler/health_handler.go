package handler

import (
	"gorm.io/gorm"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthHandler provides a health check endpoint
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler initializes the health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthCheck responds with the status of the API and database
// @Summary Health Check
// @Description Returns API status and database connectivity check
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c echo.Context) error {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "unhealthy", "error": "Database connection error"})
	}

	if err := sqlDB.Ping(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "unhealthy", "error": "Database unreachable"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}
