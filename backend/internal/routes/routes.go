package routes

import (
	"github.com/PingTower/internal/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(api *echo.Group, sh *handlers.ServiceHandler) {
	api.POST("/service", sh.CreateServiceHandler)
	api.GET("/service", nil)
}
