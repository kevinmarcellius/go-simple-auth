package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/kevinmarcellius/go-simple-auth/config"
)

type HealthHandler struct {
	DB *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{DB: db}
}

func (h *HealthHandler) ReadinessCheck(c echo.Context) error {
	log.Println("Performing readiness check")
	err := config.DBHealthCheck(h.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "unhealthy"})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}

func (h *HealthHandler) LivenessCheck(c echo.Context) error {
	log.Println("Performing liveness check")
	return c.JSON(http.StatusOK, map[string]string{"status": "alive"})
}
