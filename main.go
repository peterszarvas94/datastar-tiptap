package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

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

	e := echo.New()
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return renderTemplate(c, "index", map[string]any{
			"Title": "Datastar + Tiptap",
		})
	})

	e.GET("/content-status", func(c echo.Context) error {
		content, err := loadContent(db)
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
		content, err := loadContent(db)
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
		if err := saveContent(db, payload.HTML); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]bool{"success": true})
	})

	e.Logger.Fatal(e.Start(":8080"))
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

func saveContent(db *sql.DB, html string) error {
	_, err := db.Exec(
		`INSERT INTO content (id, html) VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET html = excluded.html`,
		html,
	)
	return err
}
