package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	store := newContentStore()

	e := echo.New()
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		clientID, err := ensureClientID(c)
		if err != nil {
			return err
		}
		content, err := store.loadContent(clientID)
		if err != nil {
			return err
		}

		return templates.ExecuteTemplate(c.Response().Writer, "index", map[string]any{
			"RawPreview":      content,
			"RenderedPreview": template.HTML(content),
		})
	})

	e.GET("/content", func(c echo.Context) error {
		clientID, err := ensureClientID(c)
		if err != nil {
			return err
		}
		content, err := store.loadContent(clientID)
		if err != nil {
			return err
		}

		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		sse.PatchElements(content,
			datastar.WithSelector("#editor > div"),
			datastar.WithMode(datastar.ElementPatchModeInner),
		)

		return nil
	})

	e.PATCH("/content", func(c echo.Context) error {
		clientID, err := ensureClientID(c)
		if err != nil {
			return err
		}

		var signals SaveSignals
		err = datastar.ReadSignals(c.Request(), &signals)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid signals"})
		}

		content := trimTrailingP(signals.EditorHTML)
		store.saveContent(clientID, content)

		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		sse.PatchSignals(fmt.Appendf(nil, `{"rawPreview": "%s"}`, stripNewlines(content)))

		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))
}
