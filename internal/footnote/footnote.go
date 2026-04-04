package footnote

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Footnote interface {
	AssignFootnotesMarker(canv canvas.Canvas, pos entity.Coordinate, defaultX int, cursor VerseLineCursor, verseFootnote map[int]map[int]repository.VerseFootNotes)
	RenderVerseFootnotes(canv canvas.Canvas, y *int, footnotes map[int]map[int]repository.VerseFootNotes)
	RenderMusicFootnotes(ctx context.Context, canv canvas.Canvas, metadata *repository.HymnMetadata, y int)
	RenderTitleFootnotes(canv canvas.Canvas, y int, metadata repository.HymnData)
}

type footnoteInteractor struct {
	li lyric.Lyric
}

func New(li lyric.Lyric) Footnote {
	return &footnoteInteractor{
		li: li,
	}
}
