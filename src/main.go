package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"
)

func main() {
	_ = godotenv.Load()
	basePath := normalizeBasePath(os.Getenv("BASE_PATH"))

	store := newContentStore()

	e := echo.New()
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return renderTemplate(c, "index", map[string]any{
			"BasePath": basePath,
		})
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

		sendContentPreviewUpdates(sse, rawContent)

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

		rawContent := trimTrailingParagraph(signals.EditorHTML)

		store.saveContent(clientID, rawContent)

		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		sendContentPreviewUpdates(sse, rawContent)

		return nil
	})

	e.Logger.Fatal(e.Start(":3000"))
}
