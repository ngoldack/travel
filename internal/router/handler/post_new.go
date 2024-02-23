package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ngoldack/travel/views/pages"
)

func PostNewHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(200, "New Post", pages.PostNew())
	}
}
