package rhythm_test

import (
	"context"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/stretchr/testify/assert"
)

func TestMergeNotes(t *testing.T) {

	newNote := func(pitch string, withTies ...bool) musicxml.Note {
		result := musicxml.Note{Pitch: struct {
			Step   string `xml:"step"`
			Octave int    `xml:"octave"`
		}{
			Step: "C",
		},
			Type: musicxml.NoteLengthQuarter,
		}

		if len(withTies) > 0 && withTies[0] {
			result.Notations = &musicxml.NoteNotation{
				Tied: &musicxml.Tie{},
			}
		}

		return result
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		note     musicxml.Note
		nextNote musicxml.Note
		ts       timesig.Time
		want     float64
		want2    musicxml.Note
	}{
		{
			name:     "No tie on nextNote",
			note:     newNote("C"),
			nextNote: newNote("C"),
			ts:       timesig.Time{ /* mock implementation */ },
			want:     0,
			want2:    newNote("C"),
		},
		{
			name:     "Different pitches",
			note:     newNote("C"),
			nextNote: newNote("D"),
			ts:       timesig.Time{ /* mock implementation */ },
			want:     0,
			want2:    newNote("C"),
		},
		{
			name: "Tied and same pitch, length >= 3",
			note: musicxml.Note{Pitch: struct {
				Step   string `xml:"step"`
				Octave int    `xml:"octave"`
			}{
				Step: "C",
			},
				Type: musicxml.NoteLengthWhole,
				Notations: &musicxml.NoteNotation{
					Tied: &musicxml.Tie{},
				},
			},
			nextNote: musicxml.Note{Pitch: struct {
				Step   string `xml:"step"`
				Octave int    `xml:"octave"`
			}{
				Step: "C",
			},
				Type: musicxml.NoteLengthHalf,
				Notations: &musicxml.NoteNotation{
					Tied: &musicxml.Tie{},
				},
				Dot: []*musicxml.Dot{{}},
			},
			ts:   timesig.Time{Beat: 4, BeatType: 4},
			want: 7, // 2.0 + 2.0
			want2: musicxml.Note{
				Pitch: struct {
					Step   string `xml:"step"`
					Octave int    `xml:"octave"`
				}{
					Step: "C",
				},
				Type:      musicxml.NoteLengthWhole,
				Notations: &musicxml.NoteNotation{Tied: &musicxml.Tie{}},
			},
		},
		{
			name:     "Tied and same pitch, length < 3",
			note:     newNote("C", true),
			nextNote: newNote("C", true),
			ts:       timesig.Time{Beat: 4, BeatType: 4},
			want:     2.0, // 1.0 + 1.0
			want2: musicxml.Note{
				Pitch: struct {
					Step   string `xml:"step"`
					Octave int    `xml:"octave"`
				}{
					Step: "C",
				},
				Type:      musicxml.NoteLengthQuarter,
				Notations: &musicxml.NoteNotation{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real test, you would need to mock context.Context and timesig.Time
			// and potentially mock TransferStopSlurAndBreathmark if it's not exported/available.
			// Assuming necessary mocks/setup are in place for this test structure.
			gotLen, gotNote := rhythm.MergeNotes(context.Background(), tt.note, tt.nextNote, tt.ts)

			if gotLen != tt.want {
				t.Errorf("MergeNotes() length = %v, want %v", gotLen, tt.want)
			}

			assert.Equal(t, tt.want2, gotNote, "merged note")
		})
	}
}
