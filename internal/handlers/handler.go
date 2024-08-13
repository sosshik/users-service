package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/sosshik/users-service/internal/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *echo.Echo {
	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is healthy:)")
	})

	g := e.Group("/users")

	{
		g.POST("", h.HandleCreateUser)
		g.PUT("/:id", h.HandleUpdateUser)
		g.DELETE("/:id", h.HandleDeleteUser)
		g.GET("", h.HandleGetUsers)
		g.GET("/:id", func(c echo.Context) error { return nil })
	}

	return e
}
