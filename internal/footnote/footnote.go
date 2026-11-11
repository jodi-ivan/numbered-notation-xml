package footnote

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type TextCalc func(txt string) float64

type Footnote interface {
	AssignFootnotesMarker(ctx context.Context, canv canvas.Canvas, pos entity.Coordinate, defaultX int, cursor VerseLineCursor, verseFootnote map[int]map[int]repository.VerseFootNotes)
	RenderVerseFootnotes(ctx context.Context, canv canvas.Canvas, y *int, footnotes map[int]map[int]repository.VerseFootNotes)
	RenderMusicFootnotes(ctx context.Context, canv canvas.Canvas, metadata *repository.HymnMetadata, y int)
	RenderTitleFootnotes(ctx context.Context, canv canvas.Canvas, y int, metadata repository.HymnData)
}

type footnoteInteractor struct {
	TextWidthCalculator TextCalc
}

func New(txtCalc TextCalc) Footnote {
	return &footnoteInteractor{
		TextWidthCalculator: txtCalc,
	}
}
