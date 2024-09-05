package lyric

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type LyricPosition struct {
	Coordinate entity.Coordinate
	Lyrics     entity.Lyric
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

type VerseInfo struct {
	MarginBottom int
}

type VerseRowStyle int

const (
	VerseRowStyleSingleColumn VerseRowStyle = 12
	VerseRowStyleDualColumn   VerseRowStyle = 6
)
