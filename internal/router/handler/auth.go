package handler

import (
	"github.com/labstack/echo/v4"
)

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("user", nil)
		return c.Redirect(302, "/")
	}
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(200, "login", nil)
	}
}
