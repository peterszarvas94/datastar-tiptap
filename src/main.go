package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
)

func main() {
	store := newContentStore()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e := echo.New()
	e.Renderer = renderer
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		clientID, err := ensureClientID(c)
		if err != nil {
			return err
		}
		content, err := loadContent(store, clientID)
		if err != nil {
			return err
		}
		return renderTemplate(c, "index", map[string]any{
			"ContentPreview":  content,
			"RenderedPreview": template.HTML(content),
		})
	})

	e.GET("/content", func(c echo.Context) error {
		clientID, err := ensureClientID(c)
		if err != nil {
			return err
		}
		content, err := loadContent(store, clientID)
		if err != nil {
			return err
		}

		sse := datastar.NewSSE(c.Response().Writer, c.Request())

		sse.PatchElements(content,
			datastar.WithSelector("#editor > div"),
			datastar.WithMode(datastar.ElementPatchModeInner),
		)

		el, err := renderTemplateFragment(c, "rendered-html", map[string]any{
			"RenderedPreview": template.HTML(content),
		})
		if err != nil {
			return err
		}
		sse.PatchElements(el)

		// stipping newlines so datastar doesn't complain, beautify runs in the frontend
		sse.PatchSignals(fmt.Appendf(nil, `{"contentPreview": "%s"}`, stripNewlines(content)))

		return nil
	})

	e.PATCH("/content", func(c echo.Context) error {
		clientID, err := ensureClientID(c)
		if err != nil {
			return err
		}
		// patch
		var signals SaveSignals
		if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid signals"})
		}

		content := trimTrailingP(signals.EditorHTML)
		saveContent(store, clientID, content)

		// update ui
		sse := datastar.NewSSE(c.Response().Writer, c.Request())

		el, err := renderTemplateFragment(c, "rendered-html", map[string]any{
			"RenderedPreview": template.HTML(content),
		})
		if err != nil {
			return err
		}
		sse.PatchElements(el)

		// stipping newlines so datastar doesn't complain, beautify runs in the frontend
		sse.PatchSignals(fmt.Appendf(nil, `{"contentPreview": "%s"}`, stripNewlines((content))))

		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))
}
