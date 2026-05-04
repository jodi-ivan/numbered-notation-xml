package rhythm

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm/splitter"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Rhythm interface {
	AdjustMultiDottedRenderer(notes []*entity.NoteRenderer, x int, y int, ks keysig.KeySignature) (int, int)
	SetRhythmNotation(noteRenderer *entity.NoteRenderer, note musicxml.Note, numberedNote int)
	RenderBezier(set []SlurBezier, canv canvas.Canvas)
	RenderSlurTies(ctx context.Context, y int, canv canvas.Canvas, notes []*entity.NoteRenderer, maxXPosition float64)
	RenderBeam(ctx context.Context, y int, canv canvas.Canvas, ts timesig.TimeSignature, notes []*entity.NoteRenderer)
}

type rhythmInteractor struct {
	BeamSplitter splitter.BeamSplitter
}

func New(beamsplitter splitter.BeamSplitter) Rhythm {
	return &rhythmInteractor{
		BeamSplitter: beamsplitter,
	}
}
