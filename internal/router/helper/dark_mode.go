package helper

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetDarkMode(c echo.Context) bool {
	cookie, err := c.Cookie("dark_mode")
	if err != nil {
		return false
	}
	return cookie.Value == "true"
}

func SetDarkMode(c echo.Context, value bool) {
	c.SetCookie(&http.Cookie{
		Name:   "dark_mode",
		Path:   "/",
		Value:  strconv.FormatBool(value),
		MaxAge: 60 * 60 * 24 * 365,
	})
}
