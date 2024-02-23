package hx

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ngoldack/travel/internal/router/helper"
)

func DarkModeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		mode := helper.GetDarkMode(c)

		helper.SetDarkMode(c, !mode)

		c.Response().Header().Set("HX-Refresh", "true")
		return c.NoContent(http.StatusNoContent)
	}
}
