package rhythm

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func Test_rhythmInteractor_SetRhythmNotation(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		note         musicxml.Note
		noteRenderer *entity.NoteRenderer
		numberedNote int
		wantNote     *entity.NoteRenderer
	}{
		{
			name:         "no notation",
			note:         musicxml.Note{},
			noteRenderer: &entity.NoteRenderer{},
			wantNote:     &entity.NoteRenderer{},
		},
		{
			name: "it has ties only",
			note: musicxml.Note{
				Notations: &musicxml.NoteNotation{
					Tied: &musicxml.Tie{
						Type: musicxml.NoteSlurTypeStart,
					},
				},
			},
			numberedNote: 1,
			noteRenderer: &entity.NoteRenderer{},
			wantNote: &entity.NoteRenderer{
				Tie: &entity.Slur{
					Number: 1,
					Type:   musicxml.NoteSlurTypeStart,
				},
			},
		},
		{
			name: "it has slur only",
			note: musicxml.Note{
				Notations: &musicxml.NoteNotation{
					Slur: []musicxml.NotationSlur{
						musicxml.NotationSlur{
							Type:   musicxml.NoteSlurTypeStart,
							Number: 1,
						},
					},
				},
			},
			numberedNote: 1,
			noteRenderer: &entity.NoteRenderer{},
			wantNote: &entity.NoteRenderer{
				Slur: map[int]entity.Slur{
					1: entity.Slur{
						Number: 1,
						Type:   musicxml.NoteSlurTypeStart,
					},
				},
			},
		},
		{
			name: "it has slur and it is hoppin'",
			note: musicxml.Note{
				Notations: &musicxml.NoteNotation{
					Slur: []musicxml.NotationSlur{
						musicxml.NotationSlur{
							Type:   musicxml.NoteSlurTypeStart,
							Number: 1,
						},
						musicxml.NotationSlur{
							Type:   musicxml.NoteSlurTypeStop,
							Number: 1,
						},
					},
				},
			},
			numberedNote: 1,
			noteRenderer: &entity.NoteRenderer{},
			wantNote: &entity.NoteRenderer{
				Slur: map[int]entity.Slur{
					1: entity.Slur{
						Number: 1,
						Type:   musicxml.NoteSlurTypeHop,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ri rhythmInteractor
			ri.SetRhythmNotation(tt.noteRenderer, tt.note, tt.numberedNote)
			assert.Equal(t, tt.wantNote, tt.noteRenderer)
		})
	}
}
