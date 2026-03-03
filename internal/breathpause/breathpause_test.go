package breathpause

import (
	"context"
	"encoding/xml"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/stretchr/testify/assert"
)

func Test_breathPauseInteractor_SetAndGetBreathPauseRenderer(t *testing.T) {
	type args struct {
		noteRenderer *entity.NoteRenderer
		note         musicxml.Note
	}
	tests := []struct {
		name       string
		args       args
		want       *entity.NoteRenderer
		setNewLine bool
	}{
		{
			name: "no breath",
			args: args{
				noteRenderer: &entity.NoteRenderer{},
			},
		},
		{
			name: "has breath with no new line",
			args: args{
				note: musicxml.Note{
					Notations: &musicxml.NoteNotation{
						Articulation: &musicxml.NotationArticulation{
							BreathMark: &struct {
								Name xml.Name
							}{},
						},
					},
				},
				noteRenderer: &entity.NoteRenderer{
					MeasureNumber: 1,
				},
			},
			want: &entity.NoteRenderer{
				Articulation: &entity.Articulation{
					BreathMark: &entity.ArticulationTypesBreathMark,
				},
				MeasureNumber: 1,
				Width:         15,
			},
		},
		{
			name: "has breath with no new line",
			args: args{
				note: musicxml.Note{
					Notations: &musicxml.NoteNotation{
						Articulation: &musicxml.NotationArticulation{
							BreathMark: &struct {
								Name xml.Name
							}{},
						},
					},
				},
				noteRenderer: &entity.NoteRenderer{
					MeasureNumber: 1,
				},
			},
			want: &entity.NoteRenderer{
				Articulation: &entity.Articulation{
					BreathMark: &entity.ArticulationTypesBreathMark,
				},
				MeasureNumber: 1,
				Width:         15,
			},
		},
		{
			name: "has breath with new line",
			args: args{
				note: musicxml.Note{
					Notations: &musicxml.NoteNotation{
						Articulation: &musicxml.NotationArticulation{
							BreathMark: &struct {
								Name xml.Name
							}{},
						},
					},
				},
				noteRenderer: &entity.NoteRenderer{
					MeasureNumber: 1,
					IsNewLine:     true,
				},
			},
			setNewLine: true,
			want: &entity.NoteRenderer{
				Articulation: &entity.Articulation{
					BreathMark: &entity.ArticulationTypesBreathMark,
				},
				MeasureNumber: 1,
				Width:         15,
				IsNewLine:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpi := breathPauseInteractor{}
			if got := bpi.SetAndGetBreathPauseRenderer(tt.args.noteRenderer, tt.args.note); !assert.Equal(t, tt.want, got) {
				t.Errorf("breathPauseInteractor.SetAndGetBreathPauseRenderer() = %v, want %v", got, tt.want)
			}
			if tt.setNewLine {
				if !assert.False(t, tt.args.noteRenderer.IsNewLine) {
					t.Errorf("breathPauseInteractor.SetAndGetBreathPauseRenderer() failed to set the new line to false")
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		if got := New(); !assert.NotNil(t, got) {
			t.Fail()
		}
	})
}

func TestAdjustBreathmarkBeamCont(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		note *entity.NoteRenderer
		prev *entity.NoteRenderer
		next *entity.NoteRenderer

		wantNote *entity.NoteRenderer
	}{
		{
			name:     "alone and not new line",
			note:     &entity.NoteRenderer{},
			wantNote: &entity.NoteRenderer{},
		},
		{
			name: "alone and new line and not beam",
			note: &entity.NoteRenderer{
				IsNewLine: true,
			},
			prev: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{},
			},
			wantNote: &entity.NoteRenderer{
				IsNewLine: true,
				Beam:      map[int]entity.Beam{},
			},
		},
		{
			name: "alone and new line and has beam",
			note: &entity.NoteRenderer{
				IsNewLine: true,
			},
			prev: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{
					1: entity.Beam{
						Number: 1,
						Type:   musicxml.NoteBeamTypeEnd,
					},
				},
			},
			wantNote: &entity.NoteRenderer{
				IsNewLine: true,
				Beam: map[int]entity.Beam{
					1: entity.Beam{
						Number: 1,
						Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
					},
				},
			},
		},
		{
			name: "has next but prev does not have beam",
			note: &entity.NoteRenderer{},
			prev: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{},
			},
			wantNote: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{},
			},
			next: &entity.NoteRenderer{},
		},
		{
			name: "has next but prev does have beam but next does not",
			note: &entity.NoteRenderer{},
			prev: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{
					1: entity.Beam{
						Number: 1,
						Type:   musicxml.NoteBeamTypeEnd,
					},
				},
			},
			wantNote: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{},
			},
			next: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{},
			},
		},
		{
			name: "has next, prev and next do have beam",
			note: &entity.NoteRenderer{},
			prev: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{
					1: entity.Beam{
						Number: 1,
						Type:   musicxml.NoteBeamTypeContinue,
					},
				},
			},
			wantNote: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{
					1: entity.Beam{
						Number: 1,
						Type:   musicxml.NoteBeam_INTERNAL_TypeAdditional,
					},
				},
			},
			next: &entity.NoteRenderer{
				Beam: map[int]entity.Beam{
					1: entity.Beam{
						Number: 1,
						Type:   musicxml.NoteBeamTypeEnd,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AdjustBreathmarkBeamCont(context.Background(), tt.note, tt.prev, tt.next)
			assert.Equal(t, tt.wantNote.Beam, tt.note.Beam)
		})
	}
}

func TestIsBreathMark(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		note *entity.NoteRenderer
		want bool
	}{
		{
			name: "nope",
			note: &entity.NoteRenderer{},
			want: false,
		},
		{
			name: "yes",
			note: &entity.NoteRenderer{
				Articulation: &entity.Articulation{
					BreathMark: &entity.ArticulationTypesBreathMark,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBreathMark(tt.note)

			assert.Equal(t, tt.want, got)
		})
	}
}
