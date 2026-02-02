package main

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func renderTemplateFragment(c echo.Context, templateName string, data any) (string, error) {
	renderer, ok := c.Echo().Renderer.(*TemplateRenderer)
	if !ok || renderer == nil {
		return "", echo.ErrInternalServerError
	}
	var buf bytes.Buffer
	err := renderer.templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderTemplate(c echo.Context, templateName string, data any) error {
	return c.Render(http.StatusOK, templateName, data)
}
