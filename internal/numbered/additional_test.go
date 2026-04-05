package numbered

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func Test_numberedInteractor_RendererFromAdditional(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		note        musicxml.Note
		header      *entity.NoteRenderer
		additionals []NoteLength
		want        []*entity.NoteRenderer
	}{
		{
			name:   "One note - quarter",
			header: &entity.NoteRenderer{},
			additionals: []NoteLength{
				{Type: musicxml.NoteLengthQuarter},
			},
			want: []*entity.NoteRenderer{
				{},
			},
		},
		{
			name: "One note - eighth",
			header: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{},
			},
			additionals: []NoteLength{
				{Type: musicxml.NoteLengthEighth},
			},
			want: []*entity.NoteRenderer{
				{
					Beam: map[int]entity.Beam{
						1: entity.Beam{
							Number: 1,
							Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
						},
					},
				},
			},
		},
		{
			name: "One note - 16th",
			header: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{},
			},
			additionals: []NoteLength{
				{Type: musicxml.NoteLength16th},
			},
			want: []*entity.NoteRenderer{
				{
					Beam: map[int]entity.Beam{
						1: entity.Beam{
							Number: 1,
							Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
						},
						2: entity.Beam{
							Number: 2,
							Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
						},
					},
				},
			},
		},
		{
			name: "two notes quater-quater",
			header: &entity.NoteRenderer{
				PositionY:     100,
				MeasureNumber: 1,
				Beam:          map[int]entity.Beam{},
			},
			additionals: []NoteLength{
				{Type: musicxml.NoteLengthQuarter},
				{Type: musicxml.NoteLengthQuarter, IsDotted: true},
			},
			want: []*entity.NoteRenderer{
				{
					PositionY:     100,
					MeasureNumber: 1,
					Beam:          map[int]entity.Beam{},
				},
				{
					PositionY:     100,
					MeasureNumber: 1,
					Beam:          map[int]entity.Beam{},
					NoteLength:    musicxml.NoteLengthQuarter,
					IsDotted:      true,
					Width:         15,
				},
			},
		},
		{
			name: "two notes quater-eighth",
			header: &entity.NoteRenderer{
				PositionY:     100,
				MeasureNumber: 1,
				Beam:          map[int]entity.Beam{},
			},
			additionals: []NoteLength{
				{Type: musicxml.NoteLengthQuarter},
				{Type: musicxml.NoteLengthEighth, IsDotted: true},
			},
			want: []*entity.NoteRenderer{
				{
					PositionY:     100,
					MeasureNumber: 1,
					Beam:          map[int]entity.Beam{},
				},
				{
					PositionY:     100,
					MeasureNumber: 1,
					Beam: map[int]entity.Beam{
						1: entity.Beam{
							Number: 1,
							Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
						},
					},
					NoteLength: musicxml.NoteLengthEighth,
					IsDotted:   true,
					Width:      15,
				},
			},
		},
		{
			name: "two notes quater-16th",
			header: &entity.NoteRenderer{
				PositionY:     100,
				MeasureNumber: 1,
				Beam:          map[int]entity.Beam{},
				IsNewLine:     true,
			},
			additionals: []NoteLength{
				{Type: musicxml.NoteLengthQuarter},
				{Type: musicxml.NoteLength16th, IsDotted: true},
			},
			want: []*entity.NoteRenderer{
				{
					PositionY:     100,
					MeasureNumber: 1,
					Beam:          map[int]entity.Beam{},
				},
				{
					PositionY:     100,
					MeasureNumber: 1,
					Beam: map[int]entity.Beam{
						1: entity.Beam{
							Number: 1,
							Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
						},
						2: entity.Beam{
							Number: 2,
							Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
						},
					},
					NoteLength: musicxml.NoteLength16th,
					IsDotted:   true,
					Width:      15,
					IsNewLine:  true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ni numberedInteractor
			got := ni.RendererFromAdditional(tt.note, tt.header, tt.additionals)
			assert.Equal(t, tt.want, got, "RendererFromAdditional()")

		})
	}
}
