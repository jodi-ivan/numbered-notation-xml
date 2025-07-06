package breathpause

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

type BreathPause interface {
	// XXX: get the proper name for this
	SetAndGetBreathPauseRenderer(noteRenderer *entity.NoteRenderer, note musicxml.Note) *entity.NoteRenderer
}

type breathPauseInteractor struct{}

func New() BreathPause {
	return &breathPauseInteractor{}
}

func (bpi *breathPauseInteractor) SetAndGetBreathPauseRenderer(noteRenderer *entity.NoteRenderer, note musicxml.Note) *entity.NoteRenderer {
	hasBreathMark := note.Notations != nil &&
		note.Notations.Articulation != nil &&
		note.Notations.Articulation.BreathMark != nil

	if !hasBreathMark {
		return nil
	}

	result := &entity.NoteRenderer{
		Articulation: &entity.Articulation{
			BreathMark: &entity.ArticulationTypesBreathMark,
		},
		MeasureNumber: noteRenderer.MeasureNumber,
		Width:         constant.LOWERCASE_LENGTH,

		// move the new line indicator to this
		IsNewLine: noteRenderer.IsNewLine,
	}

	if noteRenderer.IsNewLine {
		// remove the new line, since it is transferred to the breath mark
		noteRenderer.IsNewLine = false
	}

	return result
}
