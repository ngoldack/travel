package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/ngoldack/travel/internal/router/helper"
	"github.com/ngoldack/travel/views/layouts"
	"github.com/ngoldack/travel/views/pages"
)

func GetReisen() echo.HandlerFunc {
	return func(c echo.Context) error {
		reisenC := templ.Component(pages.Reisen())
		page := templ.Component(layouts.Layout(helper.GetDarkMode(c), "Reisen", reisenC))

		return page.Render(c.Request().Context(), c.Response().Writer)
	}
}
