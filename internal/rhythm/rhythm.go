package rhythm

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type Rhythm interface {
	AdjustMultiDottedRenderer(notes []*entity.NoteRenderer, x int, y int) (int, int)
	SetRhythmNotation(noteRenderer *entity.NoteRenderer, note musicxml.Note, numberedNote int)
}

type rhythmInteractor struct{}

func (ri *rhythmInteractor) AdjustMultiDottedRenderer(notes []*entity.NoteRenderer, x int, y int) (int, int) {
	return AdjustMultiDottedRenderer(notes, x, y)
}

func (ri *rhythmInteractor) SetRhythmNotation(noteRenderer *entity.NoteRenderer, note musicxml.Note, numberedNote int) {
	SetRhythmNotation(noteRenderer, note, numberedNote)
}

func New() Rhythm {
	return &rhythmInteractor{}
}
