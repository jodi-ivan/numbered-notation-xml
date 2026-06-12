package entity

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
)

type LyricStylePart struct {
	Text      string `json:"text"`
	Underline bool   `json:"underline"`
}

type LyricPartVerse struct {
	Text      string                 `json:"text"`
	Type      musicxml.LyricSyllabic `json:"type"`
	Combine   bool                   `json:"combine"`
	Breakdown []LyricStylePart       `json:"breakdown"`
	Offset    int                    `json:"offset"`
}

type LyricWordVerse struct {
	Word      string           `json:"word"`
	Breakdown []LyricPartVerse `json:"breakdown"`
	Dash      bool             `json:"dash"`
}

type HymnMetaData struct {
	*repository.HymnMetadata
	ParsedVerse map[int][][]LyricWordVerse
}
