package middleware

import (
	"net/http/httptest"

	"github.com/betacraft/yaag/middleware"
	"github.com/betacraft/yaag/yaag"
	"github.com/betacraft/yaag/yaag/models"
	"github.com/labstack/echo"
)

func Yaag() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			if !yaag.IsOn() {
				return next(c)
			}

			apiCall := models.ApiCall{}
			writer := httptest.NewRecorder()
			oldWriter := c.Response().Writer
			c.Response().Writer = writer
			middleware.Before(&apiCall, c.Request())
			err := next(c)
			c.Response().Writer = oldWriter
			middleware.After(&apiCall, writer, c.Response(), c.Request())
			return err
		})
	})
}
