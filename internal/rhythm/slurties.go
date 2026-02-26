package rhythm

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

func (ri *rhythmInteractor) SetRhythmNotation(noteRenderer *entity.NoteRenderer, note musicxml.Note, numberedNote int) {
	if note.Notations != nil {

		for i, slur := range note.Notations.Slur {
			if i == 0 {
				noteRenderer.Slur = map[int]entity.Slur{}
			}

			_, existing := noteRenderer.Slur[slur.Number]
			if !existing {
				noteRenderer.Slur[slur.Number] = entity.Slur{
					Number:   slur.Number,
					Type:     slur.Type,
					LineType: slur.LineType,
				}
			} else {
				noteRenderer.Slur[slur.Number] = entity.Slur{
					Number: slur.Number,
					Type:   musicxml.NoteSlurTypeHop,
				}
			}

		}

		if note.Notations.Tied != nil {
			noteRenderer.Tie = &entity.Slur{
				Number: numberedNote,
				Type:   note.Notations.Tied.Type,
			}
		}

		noteRenderer.Tuplet = note.Notations.Tuplet
	}
}

func TransferStopSlurAndBreathmark(from, to musicxml.Note) musicxml.Note {
	for _, v := range from.Notations.Slur {
		if v.Type == musicxml.NoteSlurTypeStop {
			to.Notations.Slur = from.Notations.Slur
			break
		}
	}
	if from.Notations != nil && from.Notations.Articulation != nil {
		if to.Notations != nil {
			to.Notations.Articulation = from.Notations.Articulation
		} else {
			to.Notations = &musicxml.NoteNotation{
				Articulation: from.Notations.Articulation,
			}
		}
	}

	return to
}

func HasTies(note musicxml.Note) bool {
	return note.Notations != nil && note.Notations.Tied != nil
}

func MergeNote(ctx context.Context, current, next musicxml.Note, ts timesig.Time) float64 {
	if !(current.Pitch.Step == next.Pitch.Step && next.Pitch.Octave == current.Pitch.Octave) {
		return 0
	}

	return ts.GetNoteLength(ctx, current) + ts.GetNoteLength(ctx, next)

}
