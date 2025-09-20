package routes

import (
	"github.com/PingTower/internal/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(api *echo.Group, sh *handlers.ServiceHandler) {
	api.POST("/services", sh.CreateServiceHandler)
	api.GET("/services", sh.GetAllServicesHandler)
	api.PUT("/services", sh.UpdateServiceHandler)
	api.DELETE("/services/:id", sh.DeleteServiceHandler)
}
