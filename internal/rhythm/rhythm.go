package rhythm

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Rhythm interface {
	AdjustMultiDottedRenderer(notes []*entity.NoteRenderer, x int, y int) (int, int)
	SetRhythmNotation(noteRenderer *entity.NoteRenderer, note musicxml.Note, numberedNote int)
	RenderBezier(set []SlurBezier, canv canvas.Canvas)
	RenderSlurTies(ctx context.Context, canv canvas.Canvas, notes []*entity.NoteRenderer, maxXPosition float64)
	RenderBeam(ctx context.Context, canv canvas.Canvas, notes []*entity.NoteRenderer)
}

type rhythmInteractor struct{}

func New() Rhythm {
	return &rhythmInteractor{}
}
