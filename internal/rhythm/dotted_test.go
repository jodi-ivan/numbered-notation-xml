package rhythm

import (
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/stretchr/testify/assert"
)

func Test_rhythmInteractor_AdjustMultiDottedRenderer(t *testing.T) {
	type args struct {
		notes []*entity.NoteRenderer
		x     int
		y     int
	}
	tests := []struct {
		name                string
		args                args
		wantX               int
		wantY               int
		wantRevisedRenderer []*entity.NoteRenderer
	}{
		{
			name: "empty renderer",
			args: args{
				notes: []*entity.NoteRenderer{},
				x:     25,
				y:     125,
			},
			wantX:               25,
			wantY:               125,
			wantRevisedRenderer: []*entity.NoteRenderer{},
		},
		{
			name: "one dotted",
			args: args{
				notes: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						IsDotted: true,
						Width:    constant.LOWERCASE_LENGTH,
					},
				},
				x: 25,
				y: 125,
			},
			wantX: 40,
			wantY: 125,
			wantRevisedRenderer: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					IsDotted:  true,
					PositionX: 20,
					PositionY: 125,
					Width:     constant.LOWERCASE_LENGTH,
				},
			},
		},
		{
			name: "two dotted",
			args: args{
				notes: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						IsDotted: true,
					},
					&entity.NoteRenderer{
						IsDotted: true,
					},
				},
				x: 25,
				y: 125,
			},
			wantX: 40,
			wantY: 125,
			wantRevisedRenderer: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					IsDotted:  true,
					PositionX: 20,
					PositionY: 125,
				},
				&entity.NoteRenderer{
					IsDotted:      true,
					PositionX:     40,
					PositionY:     125,
					IndexPosition: 1,
				},
			},
		},
		{
			name: "two dotted - first note width is taken from lyric",
			args: args{
				notes: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						IsDotted:               true,
						Width:                  50,
						IsLengthTakenFromLyric: true,
					},
					&entity.NoteRenderer{
						IsDotted: true,
					},
				},
				x: 25,
				y: 125,
			},
			wantX: 90,
			wantY: 125,
			wantRevisedRenderer: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					IsDotted:               true,
					Width:                  50,
					PositionX:              20,
					PositionY:              125,
					IsLengthTakenFromLyric: true,
				},
				&entity.NoteRenderer{
					IsDotted:      true,
					PositionX:     40,
					PositionY:     125,
					IndexPosition: 1,
				},
			},
		},
		{
			name: "two dotted and need to move the next note after the second",
			args: args{
				notes: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 25,
						IsDotted:  true,
						Width:     constant.LOWERCASE_LENGTH,
					},
					&entity.NoteRenderer{
						PositionX: 40,
						Width:     constant.LOWERCASE_LENGTH,
						IsDotted:  true,
					},
					&entity.NoteRenderer{
						Width:     constant.LOWERCASE_LENGTH,
						PositionX: 55,
						IsDotted:  false,
					},
				},
				x: 25,
				y: 125,
			},
			wantX: 85,
			wantY: 125,
			wantRevisedRenderer: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					IsDotted:  true,
					Width:     constant.LOWERCASE_LENGTH,
					PositionX: 20,
					PositionY: 125,
				},
				&entity.NoteRenderer{
					IsDotted:      true,
					Width:         constant.LOWERCASE_LENGTH,
					PositionX:     40,
					PositionY:     125,
					IndexPosition: 1,
				},
				&entity.NoteRenderer{
					IsDotted:      false,
					Width:         constant.LOWERCASE_LENGTH,
					PositionX:     70,
					PositionY:     125,
					IndexPosition: 2,
				},
			},
		},
		{
			name: "two dotted and need to move the next note is breath mark",
			args: args{
				notes: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 25,
						IsDotted:  true,
						Width:     constant.LOWERCASE_LENGTH,
					},
					&entity.NoteRenderer{
						PositionX: 40,
						Width:     constant.LOWERCASE_LENGTH,
						IsDotted:  true,
					},
					&entity.NoteRenderer{
						Width:     constant.LOWERCASE_LENGTH,
						PositionX: 55,
						IsDotted:  false,
						Articulation: &entity.Articulation{
							BreathMark: &entity.ArticulationTypesBreathMark,
						},
					},
				},
				x: 25,
				y: 125,
			},
			wantX: 70,
			wantY: 125,
			wantRevisedRenderer: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					IsDotted:  true,
					Width:     constant.LOWERCASE_LENGTH,
					PositionX: 20,
					PositionY: 125,
				},
				&entity.NoteRenderer{
					IsDotted:      true,
					Width:         constant.LOWERCASE_LENGTH,
					PositionX:     40,
					PositionY:     125,
					IndexPosition: 1,
				},
				&entity.NoteRenderer{
					IsDotted:      false,
					Width:         constant.LOWERCASE_LENGTH,
					PositionX:     55,
					PositionY:     125,
					IndexPosition: 2,
					Articulation: &entity.Articulation{
						BreathMark: &entity.ArticulationTypesBreathMark,
					},
				},
			},
		},
		{
			name: "two dotted and need to move the next note is new line",
			args: args{
				notes: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 25,
						IsDotted:  true,
						Width:     constant.LOWERCASE_LENGTH,
					},
					&entity.NoteRenderer{
						PositionX: 40,
						Width:     constant.LOWERCASE_LENGTH,
						IsDotted:  true,
					},
					&entity.NoteRenderer{
						Width:     constant.LOWERCASE_LENGTH,
						PositionX: 55,
						IsDotted:  false,
						IsNewLine: true,
					},
				},
				x: 25,
				y: 125,
			},
			wantX: 50,
			wantY: 125,
			wantRevisedRenderer: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					IsDotted:  true,
					Width:     constant.LOWERCASE_LENGTH,
					PositionX: 20,
					PositionY: 125,
				},
				&entity.NoteRenderer{
					IsDotted:      true,
					Width:         constant.LOWERCASE_LENGTH,
					PositionX:     40,
					PositionY:     125,
					IndexPosition: 1,
				},
				&entity.NoteRenderer{
					IsDotted:      false,
					Width:         constant.LOWERCASE_LENGTH,
					PositionX:     70,
					PositionY:     125,
					IndexPosition: 2,
					IsNewLine:     true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ri := &rhythmInteractor{}
			gotX, gotY := ri.AdjustMultiDottedRenderer(tt.args.notes, tt.args.x, tt.args.y)
			if gotX != tt.wantX {
				t.Errorf("rhythmInteractor.AdjustMultiDottedRenderer() got X = %v, want %v", gotX, tt.wantX)
			}
			if gotY != tt.wantY {
				t.Errorf("rhythmInteractor.AdjustMultiDottedRenderer() got Y = %v, want %v", gotY, tt.wantY)
			}
			assert.Equal(t, tt.wantRevisedRenderer, tt.args.notes)
		})
	}
}
