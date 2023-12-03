package lyric

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type LyricPosition struct {
	Coordinate entity.Coordinate
	Lyrics     entity.Lyric
}

type LyricPartVerse struct {
	Text    string                 `json:"text"`
	Type    musicxml.LyricSyllabic `json:"type"`
	Combine bool                   `json:"combine"`
	Dash    bool                   `json:"dash"`
}

type LyricWordVerse struct {
	Word      string           `json:"word"`
	Breakdown []LyricPartVerse `json:"breakdown"`
}
