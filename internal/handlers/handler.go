package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) InitRoutes() *echo.Echo {

	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is healthy:)")
	})
	return e
}
