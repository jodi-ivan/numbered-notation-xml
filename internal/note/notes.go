package note

import (
	"context"

	"github.com/google/uuid"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/moveabledo"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

func ConvertToRenderer(ctx context.Context, x, y, measureIdx, pos int, measure musicxml.Measure, ks keysig.KeySignature, ts timesig.TimeSignature) *entity.NoteRenderer {
	note := measure.Notes[measureIdx]
	n, octave, strikethrough := moveabledo.GetNumberedNotation(ks.GetKeyOnMeasure(ctx, measure.Number), note)
	noteLength := ts.GetNoteLength(ctx, measure.Number, note)

	renderer := &entity.NoteRenderer{
		UUID:      uuid.New().String(),
		PositionX: x, PositionY: int(y),
		Note: n, NoteLength: note.Type, Octave: octave,
		Strikethrough: strikethrough, NoteValue: noteLength,
		IsRest: (note.Rest != nil),

		Beam:          map[int]entity.Beam{},
		IsNewLine:     measure.NewLineIndex[measureIdx],
		MeasureNumber: measure.Number,

		AbsoluteNote: note.Pitch.Step, AbsoluteOctave: note.Pitch.Octave,
		AbsoluteAccidental: note.Accidental,
		TimeModifications:  note.TimeModification,

		LeadingHeader: measure.PrefixHeader[measureIdx],
		// IndexPosition: pos + note.IndexPosition + data.IndexStart,
	}

	if note.Notations != nil && note.Notations.Fermata != nil {
		renderer.Fermata = note.Notations.Fermata
	}

	return renderer
}
