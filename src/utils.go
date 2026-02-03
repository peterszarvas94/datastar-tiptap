package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func trimTrailingParagraph(html string) string {
	return strings.TrimSuffix(html, "<p></p>")
}

func stripNewlines(html string) string {
	return strings.ReplaceAll(html, "\n", "")
}

func getTemplateFragment(templateName string, data any) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderTemplate(c echo.Context, templateName string, data any) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return templates.ExecuteTemplate(c.Response().Writer, templateName, data)
}

func sendContentPreviewUpdates(sse *datastar.ServerSentEventGenerator, content string) error {
	// raw
	sse.PatchSignals(fmt.Appendf(nil, `{"rawPreview": "%s"}`, stripNewlines(content)))

	// rendered
	renderedPreview, err := getTemplateFragment("rendered-preview", map[string]any{
		"RenderedPreview": template.HTML(content),
	})
	if err != nil {
		return err
	}
	sse.PatchElements(renderedPreview)

	return nil
}
