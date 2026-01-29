package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type ContentPayload struct {
	HTML string `json:"html"`
}

const contentPath = "data/content.html"

func main() {
	e := echo.New()
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Static("/assets", "static")

	e.GET("/", func(c echo.Context) error {
		return renderTemplate(c, "index", map[string]any{
			"Title": "Datastar + Tiptap",
		})
	})

	e.GET("/content-status", func(c echo.Context) error {
		content, err := loadContent()
		if err != nil {
			return err
		}
		status := fmt.Sprintf(
			`<div id="content-status">Saved content length: %d characters</div>`,
			len(content),
		)
		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		time.Sleep(500 * time.Millisecond)
		sse.PatchElements(status)
		return nil
	})

	e.GET("/api/content", func(c echo.Context) error {
		content, err := loadContent()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{"html": content})
	})

	e.PATCH("/api/content", func(c echo.Context) error {
		var payload ContentPayload
		if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if err := saveContent(payload.HTML); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]bool{"success": true})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func renderTemplate(c echo.Context, templateName string, data any) error {
	return c.Render(http.StatusOK, templateName, data)
}

func loadContent() (string, error) {
	content, err := os.ReadFile(contentPath)
	if err == nil {
		return string(content), nil
	}
	if os.IsNotExist(err) {
		return "", nil
	}
	return "", err
}

func saveContent(html string) error {
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}
	return os.WriteFile(contentPath, []byte(html), 0644)
}
