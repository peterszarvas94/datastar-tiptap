package main

import "strings"

func trimTrailingP(html string) string {

	return strings.TrimSuffix(html, "<p></p>")
}

func stripNewlines(html string) string {
	return strings.ReplaceAll(html, "\n", "")
}
