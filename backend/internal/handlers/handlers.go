package handlers

import (
	"net/http"

	"github.com/PingTower/internal/entities"
	"github.com/PingTower/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

	if err := h.manager.CreateService(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "service created"})
}

func (h *ServiceHandler) GetAllServicesHandler(c echo.Context) error {
	services, err := h.manager.GetAllServices(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, services)
}

func (h *ServiceHandler) UpdateServiceHandler(c echo.Context) error {
	var req entities.UpdateServiceRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.manager.UpdateService(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "service updated"})
}

func (h *ServiceHandler) DeleteServiceHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid UUID"})
	}

	if err := h.manager.DeleteService(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "service deleted"})
}
