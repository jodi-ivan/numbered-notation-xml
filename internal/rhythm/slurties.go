package rhythm

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

func SetRhythmNotation(noteRenderer *entity.NoteRenderer, note musicxml.Note, numberedNote int) {
	if note.Notations != nil {

		for i, slur := range note.Notations.Slur {
			if i == 0 {
				noteRenderer.Slur = map[int]entity.Slur{}
			}

			_, existing := noteRenderer.Slur[slur.Number]
			if !existing {
				noteRenderer.Slur[slur.Number] = entity.Slur{
					Number: slur.Number,
					Type:   slur.Type,
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
