package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
	_ "modernc.org/sqlite"
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

type SaveSignals struct {
	EditorHTML string `json:"editorHtml"`
}

const dbPath = "data/content.db"
const defaultContent = `<h2>Welcome to the editor</h2>
<p>This content is stored in SQLite.</p>
<blockquote><p>Edit this and click Save.</p></blockquote>`

func main() {
	db, err := openDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e := echo.New()
	e.Renderer = renderer
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		content, err := loadContent(db)
		if err != nil {
			return err
		}
		return renderTemplate(c, "index", map[string]any{
			"ContentPreview":  content,
			"RenderedPreview": template.HTML(content),
		})
	})

	e.GET("/content", func(c echo.Context) error {
		content, err := loadContent(db)
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
		// patch
		var signals SaveSignals
		if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid signals"})
		}

		content := trimTrailingP(signals.EditorHTML)
		if err := saveContent(db, content); err != nil {
			return err
		}

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

func openDB() (*sql.DB, error) {
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS content (id INTEGER PRIMARY KEY, html TEXT NOT NULL)`); err != nil {
		return nil, err
	}
	return db, nil
}

func loadContent(db *sql.DB) (string, error) {
	var content string
	err := db.QueryRow(`SELECT html FROM content WHERE id = 1`).Scan(&content)
	if err == nil {
		return content, nil
	}
	if err != sql.ErrNoRows {
		return "", err
	}
	if err := saveContent(db, defaultContent); err != nil {
		return "", err
	}
	return defaultContent, nil
}

func trimTrailingP(html string) string {
	return strings.TrimSuffix(html, "<p></p>")
}

func saveContent(db *sql.DB, html string) error {
	_, err := db.Exec(
		`INSERT INTO content (id, html) VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET html = excluded.html`,
		html,
	)
	return err
}

func stripNewlines(html string) string {
	return strings.ReplaceAll(html, "\n", "")
}
