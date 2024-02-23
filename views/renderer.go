package views

import (
	"io"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/ngoldack/travel/internal/router/helper"
	"github.com/ngoldack/travel/views/layouts"
)

type TemplRenderer struct{}

// Render implements echo.Renderer.
func (tr *TemplRenderer) Render(w io.Writer, n string, data interface{}, c echo.Context) error {
	ctx := c.Request().Context()

	content := templ.Component(data.(templ.Component))
	page := layouts.Layout(helper.GetDarkMode(c), n, content)

	return page.Render(ctx, w)
}

var _ echo.Renderer = (*TemplRenderer)(nil)
