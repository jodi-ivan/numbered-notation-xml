package numbered

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

type renderedNote struct {
	IsDotted bool
	Type     musicxml.NoteLength
}

func RenderLengthNote(ctx context.Context, ts timesig.TimeSignature, note musicxml.Note, noteLength float64) []renderedNote {
	return nil

}
