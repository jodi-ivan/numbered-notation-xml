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
	ubermock "go.uber.org/mock/gomock"
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
	uberCtrl := ubermock.NewController(t)
	defer uberCtrl.Finish()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv        func(*gomock.Controller) *canvas.MockCanvasTestify
		lyricMock   func(*gomock.Controller) *lyric.MockLyric
		barlineMock func(*ubermock.Controller) *barline.MockBarline

		measure          []*entity.NoteRenderer
		y                int
		rightAlignOffset int
	}{
		{
			name: "everything went fine",
			canv: func(c *gomock.Controller) *canvas.MockCanvasTestify {
				canv := canvas.NewMockCanvasTestify(t)
				canv.EXPECT().Text(55, 100, "2")
				canv.EXPECT().Text(60, 100, ".")
				canv.EXPECT().TextUnescaped(72.0, 90.0, ",")
				canv.EXPECT().Text(70, 100, "1")

				canv.EXPECT().Circle(59, 72, 6, []string{`stroke="black"`, `fill="none"`, `stroke-width="1.3"`})
				canv.EXPECT().Text(56, 75, "1", []string{`font-weight="600"`, `style="font-size:9.6px"`})

				return canv
			},
			barlineMock: func(c *ubermock.Controller) *barline.MockBarline {
				b := barline.NewMockBarline(c)
				b.EXPECT().RenderBarline(ubermock.Any(), ubermock.Any(), ubermock.Any(), entity.NewCoordinate(50, 100))
				return b
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
				ni.Barline = tt.barlineMock(uberCtrl)
			}
			canv := canvas.Canvas(nil)
			if tt.canv != nil {
				canv = tt.canv(ctrl)
			}
			ni.RenderNote(context.Background(), canv, tt.measure, tt.y, tt.rightAlignOffset)
			if tt.canv != nil {
				canv1 := canv.(*canvas.MockCanvasTestify)
				canv1.AssertExpectations(t)
			}
		})
	}
}

func TestNumberedInteractor_SplitNote(t *testing.T) {
	ts44 := timesig.Time{Beat: 4, BeatType: 4}
	ts34 := timesig.Time{Beat: 3, BeatType: 4}

	tests := []struct {
		name       string
		noteLength float64
		ts         timesig.Time
		flag       musicxml.NoteLength
		next       musicxml.NoteLength
		want       []NoteLength
	}{
		{
			name:       "whole note in 4/4",
			noteLength: 4,
			ts:         ts44,
			flag:       musicxml.NoteLengthWhole,
			want: []NoteLength{
				{Type: musicxml.NoteLengthWhole, IsDotted: true},
				{Type: musicxml.NoteLengthWhole, IsDotted: true},
				{Type: musicxml.NoteLengthWhole, IsDotted: true},
				{Type: musicxml.NoteLengthWhole, IsDotted: true},
			},
		},
		{
			name:       "quarter note in 4/4",
			noteLength: 1,
			ts:         ts44,
			flag:       musicxml.NoteLengthQuarter,
			want: []NoteLength{
				{Type: musicxml.NoteLengthQuarter, IsDotted: true},
			},
		},
		{
			name:       "two quarters in 4/4",
			noteLength: 2,
			ts:         ts44,
			flag:       musicxml.NoteLengthQuarter,
			want: []NoteLength{
				{Type: musicxml.NoteLengthQuarter, IsDotted: true},
				{Type: musicxml.NoteLengthQuarter, IsDotted: true},
			},
		},
		{
			name:       "eighth note in 4/4",
			noteLength: 0.5,
			ts:         ts44,
			flag:       musicxml.NoteLengthEighth,
			want: []NoteLength{
				{Type: musicxml.NoteLengthEighth, IsDotted: true},
			},
		},
		{
			// unit clamped to Quarter; not merge (unit*2=2 != 1)
			// fullNotes=floor(2/1)=2, remaining=0 → two dotted Half notes
			name:       "half flag clamped in 3/4 - two units",
			noteLength: 2,
			ts:         ts34,
			flag:       musicxml.NoteLengthHalf,
			want: []NoteLength{
				{Type: musicxml.NoteLengthHalf, IsDotted: true},
				{Type: musicxml.NoteLengthHalf, IsDotted: true},
			},
		},
		{
			// unit clamped to Quarter; fullNotes=floor(1.5/1)=1, remaining=0.5
			// dotValue=0.5, flags[0.5]=Eighth exists, remaining>=0.5 → append dotted Eighth
			name:       "half flag clamped in 3/4 - 1.5 units",
			noteLength: 1.5,
			ts:         ts34,
			flag:       musicxml.NoteLengthHalf,
			want: []NoteLength{
				{Type: musicxml.NoteLengthHalf, IsDotted: true},
				{Type: musicxml.NoteLengthEighth, IsDotted: true},
			},
		},
		{
			// kj-064: 2.5 quarters, unit=1, fullNotes=2, remaining=0.5
			// → [dotted Q, dotted Q, dotted Eighth]
			name:       "kj-064 two and a half quarters in 4/4",
			noteLength: 2.5,
			ts:         ts44,
			flag:       musicxml.NoteLengthQuarter,
			want: []NoteLength{
				{Type: musicxml.NoteLengthQuarter, IsDotted: true},
				{Type: musicxml.NoteLengthQuarter, IsDotted: true},
				{Type: musicxml.NoteLengthEighth, IsDotted: true},
			},
		},
		{
			// merge path: 3 eighths in 3/4
			// unit=0.5, unit*2=1 == base[flags[1]]=1 → shouldMerge
			// fullNotes=6 >= 2 → merge branch
			// results=[dotted Eighth]; fullNotes -= 0.5 → 5.5
			// append Repeat(dotted Quarter, int(5.5/2)=2)
			// noteLength = 2.5 - 2*0.5*2 = 0.5 == unit → append dotted Eighth
			name:       "merge path - 3 eighths in 3/4",
			noteLength: 3,
			ts:         ts34,
			flag:       musicxml.NoteLengthEighth,
			want: []NoteLength{
				{Type: musicxml.NoteLengthEighth, IsDotted: true},
				{Type: musicxml.NoteLengthQuarter, IsDotted: true},
				{Type: musicxml.NoteLengthQuarter, IsDotted: true},
				{Type: musicxml.NoteLengthEighth, IsDotted: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ni := &numberedInteractor{}
			got := ni.SplitNote(context.Background(), tt.noteLength, tt.ts, tt.flag, tt.next)
			if !assert.Equal(t, tt.want, got) {
				t.Logf("got:  %+v", got)
				t.Logf("want: %+v", tt.want)
			}
		})
	}
}
