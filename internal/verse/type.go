package verse

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type VerseInfo struct {
	MarginBottom int
}

type VerseRowStyle int

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

type ParsedVerse struct {
	Verse        []string
	ElisionMarks [][2]entity.Coordinate
	Position     versePosition
}

type ParsedVerseWithInfo struct {
	Verses        map[int]ParsedVerse
	IsMultiColumn bool
	MaxLineWidth  float64
	MaxRightPos   float64
	RowPositionY  map[int]int
}
