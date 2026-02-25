package breathpause

import (
	"context"

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

func AdjustBreathmarkBeamCont(ctx context.Context, note, prev, next *entity.NoteRenderer) {
	if note.IsNewLine {
		note.Beam = map[int]entity.Beam{}
		for beamNo := 1; beamNo < 4; beamNo++ {
			_, hasBeam := prev.Beam[beamNo] // previous note
			if !hasBeam {
				break
			}

			note.Beam[beamNo] = entity.Beam{
				Number: beamNo,
				Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
			}
		}
		return
	}

	if prev != nil && next != nil {
		// dont break the line if there is a breathmark
		note.Beam = map[int]entity.Beam{}
		for beamNo := 1; beamNo < 4; beamNo++ {
			_, hasBeam := prev.Beam[beamNo] // previous note
			if !hasBeam {
				break
			}

			_, hasBeam = next.Beam[beamNo] // next note
			if !hasBeam {
				break
			}

			note.Beam[beamNo] = entity.Beam{
				Number: beamNo,
				Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
			}
		}
	}

}
