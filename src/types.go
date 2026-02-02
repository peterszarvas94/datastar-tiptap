package main

type ContentPayload struct {
	HTML string `json:"html"`
}

type SaveSignals struct {
	EditorHTML string `json:"editorHtml"`
}
