package main

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
)

func main() {
	store := newContentStore()

	e := echo.New()
	e.Static("/static", "static")

	// var templates = template.Must(template.ParseGlob("templates/*.html"))
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer

	e.GET("/", func(c echo.Context) error {
		return renderTemplate(c, "index", map[string]any{})
	})

	e.GET("/content", func(c echo.Context) error {
		clientID, err := ensureClientID(c)
		if err != nil {

			return err
		}
		rawContent, err := store.loadContent(clientID)
		if err != nil {

			return err
		}

		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		sse.PatchElements(rawContent,
			datastar.WithSelector("#editor > div"),
			datastar.WithMode(datastar.ElementPatchModeInner),
		)

		updateContentPreviews(c, sse, rawContent)

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

		rawContent := trimTrailingP(signals.EditorHTML)

		store.saveContent(clientID, rawContent)

		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		updateContentPreviews(c, sse, rawContent)

		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))
}
