package main

import (
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
			"ContentPreview":  content,
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

		err = updateUI(sse, content)
		if err != nil {
			return err
		}

		return nil
	})

	e.PATCH("/content", func(c echo.Context) error {
		clientID, err := ensureClientID(c)
		if err != nil {
			return err
		}
		// patch
		var signals SaveSignals
		err = datastar.ReadSignals(c.Request(), &signals)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid signals"})
		}

		content := trimTrailingP(signals.EditorHTML)
		store.saveContent(clientID, content)

		// update ui
		sse := datastar.NewSSE(c.Response().Writer, c.Request())

		err = updateUI(sse, content)
		if err != nil {
			return err
		}

		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))
}
