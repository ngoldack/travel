package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/ngoldack/travel/internal/router/helper"
	"github.com/ngoldack/travel/views/layouts"
	"github.com/ngoldack/travel/views/pages"
)

func GetIndex() echo.HandlerFunc {
	return func(c echo.Context) error {
		indexC := templ.Component(pages.Index())
		err := templ.Component(layouts.Layout(helper.GetDarkMode(c), "Home", indexC)).Render(c.Request().Context(), c.Response().Writer)
		if err != nil {
			return err
		}
		return nil
	}
}
