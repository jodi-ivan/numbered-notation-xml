package verse

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
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
