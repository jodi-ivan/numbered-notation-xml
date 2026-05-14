package lyric

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
)

type LyricPosition struct {
	Coordinate entity.Coordinate
	Lyrics     entity.Lyric
	TotalLyric int
}

type VerseInfo struct {
	MarginBottom int
	HasLyric     bool
}
