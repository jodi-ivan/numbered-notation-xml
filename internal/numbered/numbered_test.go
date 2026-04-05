package numbered

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
)

func Test_numberedInteractor_GetLengthNote(t *testing.T) {
	type args struct {
		ts         timesig.TimeSignature
		measure    int
		noteLength float64
	}
	tests := []struct {
		name string
		args args
		want []NoteLength
	}{
		{
			name: "quarter .25 beat",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 0.25,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLength16th,
				},
			},
		},
		{
			name: "quarter .5 beat",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 0.5,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthEighth,
				},
			},
		},
		{
			name: "quarter .75 beat",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 0.75,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthEighth,
				},
				NoteLength{
					IsDotted: true,
					Type:     musicxml.NoteLength16th,
				},
			},
		},
		{
			name: "quarter 1 beat",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 1,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthQuarter,
				},
			},
		},
		{
			name: "quarter 1.5 beat",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 1.5,
			},
			want: []NoteLength{ /// 1 ^.
				NoteLength{
					Type: musicxml.NoteLengthQuarter,
				},
				NoteLength{
					Type:     musicxml.NoteLengthEighth,
					IsDotted: true,
				},
			},
		},
		{
			name: "quarter 2 beats",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 2,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthQuarter,
				},
				NoteLength{
					IsDotted: true,
					Type:     musicxml.NoteLengthQuarter,
				},
			},
		},
		{
			name: "quarter 2.5 beat",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 2.5,
			},
			want: []NoteLength{ /// 1 ^.
				NoteLength{
					Type: musicxml.NoteLengthQuarter,
				},
				NoteLength{
					Type:     musicxml.NoteLengthQuarter,
					IsDotted: true,
				},
				NoteLength{
					Type:     musicxml.NoteLengthEighth,
					IsDotted: true,
				},
			},
		},
		{
			name: "quarter 3 beats",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 4,
						},
					},
				},
				measure:    1,
				noteLength: 3,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthQuarter,
				},
				NoteLength{
					Type:     musicxml.NoteLengthQuarter,
					IsDotted: true,
				},
				NoteLength{
					Type:     musicxml.NoteLengthQuarter,
					IsDotted: true,
				},
			},
		},
		{
			name: "eight 1 beats",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 8,
						},
					},
				},
				measure:    1,
				noteLength: 1,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthEighth,
				},
			},
		},
		{
			name: "eight 0.5 beats",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 8,
						},
					},
				},
				measure:    1,
				noteLength: 0.5,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLength16th,
				},
			},
		},
		{
			name: "eight 0.25 beats",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 8,
						},
					},
				},
				measure:    1,
				noteLength: 0.25,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLength32nd,
				},
			},
		},
		{
			name: "eight 1.5 beats",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 8,
						},
					},
				},
				measure:    1,
				noteLength: 1.5,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthEighth,
				},
				NoteLength{
					IsDotted: true,
					Type:     musicxml.NoteLength16th,
				},
			},
		},
		{
			name: "eight 2 beats",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 8,
						},
					},
				},
				measure:    1,
				noteLength: 2,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthEighth,
				},
				NoteLength{
					IsDotted: true,
					Type:     musicxml.NoteLengthEighth,
				},
			},
		},
		{
			name: "eight 2.25 beats",
			args: args{
				ts: timesig.TimeSignature{
					Signatures: []timesig.Time{
						timesig.Time{
							BeatType: 8,
						},
					},
				},
				measure:    1,
				noteLength: 2.25,
			},
			want: []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthEighth,
				},
				NoteLength{
					IsDotted: true,
					Type:     musicxml.NoteLengthEighth,
				},
				NoteLength{
					IsDotted: true,
					Type:     musicxml.NoteLength16th,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ni := numberedInteractor{}
			if got := ni.GetLengthNote(context.Background(), tt.args.ts, tt.args.measure, tt.args.noteLength); !assert.Equal(t, tt.want, got) {
				t.Errorf("numberedInteractor.GetLengthNote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_numberedInteractor_RenderNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv        func(*gomock.Controller) *canvas.MockCanvas
		lyricMock   func(*gomock.Controller) *lyric.MockLyric
		barlineMock func(*gomock.Controller) *barline.MockBarline

		measure          []*entity.NoteRenderer
		y                int
		rightAlignOffset int
	}{
		{
			name: "everuthing went fine",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Text(55, 100, "2")
				canv.EXPECT().Text(60, 100, ".")
				canv.EXPECT().Text(72, 90, ",")
				canv.EXPECT().Text(62, 100, "1")

				canv.EXPECT().Circle(59, 72, 6, `stroke="black"`, `fill="none"`, `stroke-width="1.3"`)
				canv.EXPECT().Text(56, 75, "1", `font-weight="600"`, `style="font-size:60%"`)

				return canv
			},
			barlineMock: func(c *gomock.Controller) *barline.MockBarline {
				b := barline.NewMockBarline(c)
				b.EXPECT().RenderBarline(gomock.Any(), gomock.Any(), gomock.Any(), entity.NewCoordinate(50, 100))
				return b
			},
			lyricMock: func(c *gomock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("2").Return(8.0)
				li.EXPECT().CalculateLyricWidth("1").Return(8.0)

				return li
			},
			measure: []*entity.NoteRenderer{
				{
					PositionX: 50,
					Barline: &musicxml.Barline{
						BarStyle: musicxml.BarLineStyleRegular,
					},
				},
				{
					PositionY:     100,
					LeadingHeader: "1",
					PositionX:     55,
					Note:          2,
				},
				{
					PositionX: 60,
					IsDotted:  true,
				},
				{
					PositionX: 65,
					Articulation: &entity.Articulation{
						BreathMark: &entity.ArticulationTypesBreathMark,
					},
				},
				{
					PositionX: 70,
					Note:      1,
				},
			},
			y: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ni numberedInteractor
			if tt.lyricMock != nil {
				ni.Lyric = tt.lyricMock(ctrl)
			}
			if tt.barlineMock != nil {
				ni.Barline = tt.barlineMock(ctrl)
			}
			canv := canvas.Canvas(nil)
			if tt.canv != nil {
				canv = tt.canv(ctrl)
			}
			ni.RenderNote(context.Background(), canv, tt.measure, tt.y, tt.rightAlignOffset)
		})
	}
}
