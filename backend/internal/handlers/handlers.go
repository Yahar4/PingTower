package handlers

import (
	"github.com/PingTower/internal/entities"
	"github.com/PingTower/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ServiceHandler struct {
	manager *service.ServiceManager
}

func NewServiceHandler(manager *service.ServiceManager) *ServiceHandler {
	return &ServiceHandler{
		manager: manager,
	}
}

func (h *ServiceHandler) CreateServiceHandler(c echo.Context) error {
	var req entities.CreateServiceRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	err := h.manager.CreateService(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "service created"})
}
