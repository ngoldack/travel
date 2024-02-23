package hx

import (
	"bytes"

	"github.com/labstack/echo/v4"
	"github.com/ngoldack/travel/views/components"
	"github.com/yuin/goldmark"
)

var gm = goldmark.New(
	goldmark.WithExtensions(),
)

func EditorPreviewRendererHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		err := c.Request().ParseForm()
		if err != nil {
			return err
		}

		editor := c.FormValue("editor")
		var preview bytes.Buffer
		err = gm.Convert([]byte(editor), &preview)
		if err != nil {
			return err
		}

		previewComponent := components.EditorPreview(preview.String())
		return previewComponent.Render(c.Request().Context(), c.Response().Writer)
	}
}
