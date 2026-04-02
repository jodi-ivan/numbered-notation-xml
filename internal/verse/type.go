package verse

import "github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"

type VerseInfo struct {
	MarginBottom int
}

type VerseRowStyle int

const (
	VerseRowStyleSingleColumn VerseRowStyle = 12
	VerseRowStyleDualColumn   VerseRowStyle = 6
)

type versePosition struct {
	Col      int
	Row      int
	RowWidth int
	Style    VerseRowStyle
}

type LyricStylePart struct {
	Text      string `json:"text"`
	Underline bool   `json:"underline"`
}

type LyricPartVerse struct {
	Text      string                 `json:"text"`
	Type      musicxml.LyricSyllabic `json:"type"`
	Combine   bool                   `json:"combine"`
	Breakdown []LyricStylePart       `json:"breakdown"`
}

type LyricWordVerse struct {
	Word      string           `json:"word"`
	Breakdown []LyricPartVerse `json:"breakdown"`
	Dash      bool             `json:"dash"`
}
