package lyric

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
)

type LyricPosition struct {
	Coordinate entity.Coordinate
	Lyrics     entity.Lyric
}

type VerseInfo struct {
	MarginBottom int
}

const (
	MAX_VERSE_IN_MUSIC          = 4
	MAX_LINE_PER_VERSE_IN_MUSIC = 2

	LINE_BETWEEN_LYRIC = 20
)
