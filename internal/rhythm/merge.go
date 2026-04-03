package rhythm

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

func MergeNotes(ctx context.Context, note, nextNote musicxml.Note, ts timesig.Time) (mergedLegth float64, mergedNote musicxml.Note) {
	if HasTies(nextNote) && note.Pitch == nextNote.Pitch {

		mergedNoteLegth := ts.GetNoteLength(ctx, note) + ts.GetNoteLength(ctx, nextNote)
		if mergedNoteLegth < 3 {
			note = TransferStopSlurAndBreathmark(nextNote, note)
		}

		return mergedNoteLegth, note
	}
	return 0, note
}
