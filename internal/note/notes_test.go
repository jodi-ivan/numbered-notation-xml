package note

import (
	"context"
	"encoding/xml"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/stretchr/testify/assert"
)

func TestConvertToRenderer(t *testing.T) {
	type args struct {
		ctx        context.Context
		x          int
		y          int
		measureIdx int
		pos        int
		measure    musicxml.Measure
		ks         keysig.KeySignature
		ts         timesig.TimeSignature
	}

	// breathmark := entity.ArticulationTypesBreathMark
	tests := []struct {
		name string
		args args
		want *entity.NoteRenderer
	}{
		{
			name: "Simple Note",
			args: args{
				ctx:        context.Background(),
				x:          100,
				y:          150,
				measureIdx: 0,
				pos:        0,
				measure: musicxml.Measure{
					Number: 1,
					Notes: []musicxml.Note{
						{
							Pitch: struct {
								Step   string `xml:"step"`
								Octave int    `xml:"octave"`
							}{Step: "C", Octave: 4},
							Accidental: musicxml.NoteAccidental("natural"),
							Type:       musicxml.NoteLengthQuarter,
						},
					},
				},
				ks: keysig.KeySignature{
					Signatures: []keysig.Key{
						{Fifth: 0},
					},
				},
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						{Beat: 4, BeatType: 4},
					},
				},
			},
			want: &entity.NoteRenderer{
				UUID:               "mock-uuid", // Updated UUID
				AbsoluteNote:       "C",
				AbsoluteOctave:     4,
				AbsoluteAccidental: musicxml.NoteAccidental("natural"),
				StemDirection:      0,
				MarginTopFromStaff: 0,
				Articulation:       nil,
				Note:               1,
				NoteLength:         musicxml.NoteLengthQuarter,
				NoteValue:          1,
				PositionX:          100,
				PositionY:          150,
				MeasureNumber:      1,
				Beam:               map[int]entity.Beam{},
			},
		},
		{
			name: "Note with Breath Mark",
			args: args{
				ctx:        context.Background(),
				x:          100,
				y:          150,
				measureIdx: 0,
				pos:        0,
				measure: musicxml.Measure{
					Number: 1,
					Notes: []musicxml.Note{
						{
							Pitch: struct {
								Step   string `xml:"step"`
								Octave int    `xml:"octave"`
							}{Step: "C", Octave: 4},
							Accidental: musicxml.NoteAccidental("natural"),
							Type:       musicxml.NoteLengthQuarter,
							Notations: &musicxml.NoteNotation{
								Articulation: &musicxml.NotationArticulation{
									BreathMark: &struct {
										Name xml.Name
									}{},
								},
							},
						},
					},
				},
				ks: keysig.KeySignature{
					Signatures: []keysig.Key{
						{Fifth: 0},
					},
				},
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						{Beat: 4, BeatType: 4},
					},
				},
			},
			want: &entity.NoteRenderer{
				UUID:               "mock-uuid", // Updated UUID
				AbsoluteNote:       "C",
				AbsoluteOctave:     4,
				AbsoluteAccidental: musicxml.NoteAccidental("natural"),
				StemDirection:      0,
				MarginTopFromStaff: 0,
				Articulation:       nil,
				Note:               1,
				NoteLength:         musicxml.NoteLengthQuarter,
				NoteValue:          1,
				PositionX:          100,
				PositionY:          150,
				MeasureNumber:      1,
				Beam:               map[int]entity.Beam{},
			},
		},
		{
			name: "Note with Sharp Accidental",
			args: args{
				ctx:        context.Background(),
				x:          100,
				y:          150,
				measureIdx: 0,
				pos:        0,
				measure: musicxml.Measure{
					Number: 1,
					Notes: []musicxml.Note{
						{
							Pitch: struct {
								Step   string `xml:"step"`
								Octave int    `xml:"octave"`
							}{Step: "C", Octave: 4},
							Accidental: musicxml.NoteAccidental("sharp"),
							Type:       musicxml.NoteLengthQuarter,
						},
					},
				},
				ks: keysig.KeySignature{
					Signatures: []keysig.Key{
						{Fifth: 0},
					},
				},
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						{Beat: 4, BeatType: 4},
					},
				},
			},
			want: &entity.NoteRenderer{
				UUID:               "mock-uuid", // Updated UUID
				AbsoluteNote:       "C",
				AbsoluteOctave:     4,
				AbsoluteAccidental: musicxml.NoteAccidental("sharp"),
				StemDirection:      0,
				MarginTopFromStaff: 0,
				Articulation:       nil,
				Note:               1,
				NoteLength:         musicxml.NoteLengthQuarter,
				NoteValue:          1,
				PositionX:          100,
				PositionY:          150,
				MeasureNumber:      1,
				Strikethrough:      true,
				Beam:               map[int]entity.Beam{},
			},
		},
		{
			name: "Note with Flat Accidental",
			args: args{
				ctx:        context.Background(),
				x:          100,
				y:          150,
				measureIdx: 0,
				pos:        0,
				measure: musicxml.Measure{
					Number: 1,
					Notes: []musicxml.Note{
						{
							Pitch: struct {
								Step   string `xml:"step"`
								Octave int    `xml:"octave"`
							}{Step: "C", Octave: 4},
							Accidental: musicxml.NoteAccidental("flat"),
							Type:       musicxml.NoteLengthQuarter,
						},
					},
				},
				ks: keysig.KeySignature{
					Signatures: []keysig.Key{
						{Fifth: 0},
					},
				},
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						{Beat: 4, BeatType: 4},
					},
				},
			},
			want: &entity.NoteRenderer{
				UUID:               "mock-uuid", // Updated UUID
				AbsoluteNote:       "C",
				AbsoluteOctave:     4,
				AbsoluteAccidental: musicxml.NoteAccidental("flat"),
				StemDirection:      0,
				MarginTopFromStaff: 0,
				Articulation:       nil,
				Note:               1,
				NoteLength:         musicxml.NoteLengthQuarter,
				NoteValue:          1,
				PositionX:          100,
				PositionY:          150,
				MeasureNumber:      1,
				Strikethrough:      true,
				Beam:               map[int]entity.Beam{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToRenderer(tt.args.ctx, tt.args.x, tt.args.y, tt.args.measureIdx, tt.args.pos, tt.args.measure, tt.args.ks, tt.args.ts)
			got.UUID = "mock-uuid"
			assert.Equal(t, tt.want, got, "ConvertToRenderer()")
		})
	}
}
