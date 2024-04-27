package breathpause

import (
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
				Width:         6,
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
				Width:         6,
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
				Width:         6,
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
