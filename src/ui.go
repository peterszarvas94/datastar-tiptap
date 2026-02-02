package main

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/starfederation/datastar-go/datastar"
)

func updateUI(sse *datastar.ServerSentEventGenerator, content string) error {
	sse.PatchElements(content,
		datastar.WithSelector("#editor > div"),
		datastar.WithMode(datastar.ElementPatchModeInner),
	)

	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "rendered-html", map[string]any{
		"RenderedPreview": template.HTML(content),
	})
	if err != nil {
		return err
	}
	sse.PatchElements(buf.String())

	// stipping newlines so datastar doesn't complain, beautify runs in the frontend
	sse.PatchSignals(fmt.Appendf(nil, `{"contentPreview": "%s"}`, stripNewlines(content)))

	return nil
}
