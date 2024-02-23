package router

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(log *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError

		if he := new(echo.HTTPError); errors.Is(err, &echo.HTTPError{}) {
			if errors.As(err, &he) {
				code = he.Code
			}
		}

		log.ErrorContext(c.Request().Context(), "error handling request", slog.Any("error", err))

		_ = c.String(code, fmt.Sprintf("Error: %s", err.Error()))
	}
}
