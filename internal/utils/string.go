package utils

import "strings"

const (
	ITALIC_OPENING = "<i>"
	ITALIC_CLOSING = "</i>"

	TSPAN_OPENING = "<tspan font-style=\"italic\">"
	TSPAN_CLOSING = "</tspan>"
)

func ReplaceItalicToSpan(italic string) string {
	italic = strings.ReplaceAll(italic, ITALIC_OPENING, TSPAN_OPENING)
	return strings.ReplaceAll(italic, ITALIC_CLOSING, TSPAN_CLOSING)
}

func ReplaceItalicToSpanWithClean(italic string) (tspan string, clean string) {
	tspan = ReplaceItalicToSpan(italic)

	clean = strings.ReplaceAll(italic, ITALIC_OPENING, "")
	return tspan, strings.ReplaceAll(clean, ITALIC_CLOSING, "")
}

func CleanSpan(span string) string {
	prefix := strings.ReplaceAll(span, TSPAN_OPENING, "")
	return strings.ReplaceAll(prefix, TSPAN_CLOSING, "")
}
