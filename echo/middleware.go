package middleware

import (
	"github.com/CoryARamirez/yaag/middleware"
	"github.com/CoryARamirez/yaag/yaag"
	"github.com/CoryARamirez/yaag/yaag/models"
	"github.com/labstack/echo"
)

func Yaag() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			if !yaag.IsOn() {
				return next(c)
			}

			apiCall := models.ApiCall{}
			writer := middleware.NewResponseRecorder(c.Response().Writer)
			c.Response().Writer = writer
			middleware.Before(&apiCall, c.Request())
			err := next(c)
			middleware.After(&apiCall, writer, c.Request())
			return err
		})
	})
}
