package main

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
)

func trimTrailingP(html string) string {

	return strings.TrimSuffix(html, "<p></p>")
}

func stripNewlines(html string) string {
	return strings.ReplaceAll(html, "\n", "")
}

func updateContentPreviews(c echo.Context, sse *datastar.ServerSentEventGenerator, content string) error {
	// raw
	sse.PatchSignals(fmt.Appendf(nil, `{"rawPreview": "%s"}`, stripNewlines(content)))

	// rendered
	renderedPreview, err := renderTemplateFragment(c, "rendered-preview", map[string]any{
		"RenderedPreview": template.HTML(content),
	})
	if err != nil {
		return err
	}
	sse.PatchElements(renderedPreview)

	return nil
}
