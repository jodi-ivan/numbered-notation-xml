package lyric

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
)

func Test_lyricInteractor_SetLyricRenderer(t *testing.T) {
	type args struct {
		noteRenderer *entity.NoteRenderer
		note         musicxml.Note
	}
	tests := []struct {
		name string
		args args
		want VerseInfo

		// since the note renderer is changing, we need to assert the changes here
		wantWidth          int
		wantTakenFromLyric bool
		wantLyric          []entity.Lyric
	}{
		{
			name: "no lyric in the renderer",
			args: args{
				noteRenderer: &entity.NoteRenderer{},
				note:         musicxml.Note{},
			},
			want:      VerseInfo{},
			wantWidth: 29,
		},
		{
			name: "lyric less than the note",
			args: args{
				noteRenderer: &entity.NoteRenderer{
					IsLengthTakenFromLyric: true,
				},
				note: musicxml.Note{
					Lyric: []musicxml.Lyric{
						{
							Number: 1,
							Text: []musicxml.LyricText{
								{Value: "a"},
							},
							Syllabic: musicxml.LyricSyllabicTypeBegin,
						},
					},
				},
			},
			want:               VerseInfo{},
			wantTakenFromLyric: false,
			wantWidth:          29,
			wantLyric: []entity.Lyric{
				{
					Text: []entity.Text{
						{Value: "a"},
					},
					Syllabic: musicxml.LyricSyllabicTypeBegin,
				},
			},
		},
		{
			name: "lyric has one note, but on 2nd notes. ",
			args: args{
				noteRenderer: &entity.NoteRenderer{},
				note: musicxml.Note{
					Lyric: []musicxml.Lyric{
						{
							Number: 2,
							Text: []musicxml.LyricText{
								{Value: "Yang"},
							},
							Syllabic: musicxml.LyricSyllabicTypeSingle,
						},
					},
				},
			},
			want:               VerseInfo{},
			wantWidth:          61,
			wantTakenFromLyric: true,
			wantLyric: []entity.Lyric{
				{},
				{
					Text: []entity.Text{
						{Value: "Yang"},
					},
					Syllabic: musicxml.LyricSyllabicTypeSingle,
				},
			},
		},
		{
			name: "lyric more than the note",
			args: args{
				noteRenderer: &entity.NoteRenderer{},
				note: musicxml.Note{
					Lyric: []musicxml.Lyric{
						{
							Number: 1,
							Text: []musicxml.LyricText{
								{Value: "Yang"},
							},
							Syllabic: musicxml.LyricSyllabicTypeSingle,
						},
					},
				},
			},
			want:               VerseInfo{},
			wantWidth:          61,
			wantTakenFromLyric: true,
			wantLyric: []entity.Lyric{
				{
					Text: []entity.Text{
						{Value: "Yang"},
					},
					Syllabic: musicxml.LyricSyllabicTypeSingle,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := &lyricInteractor{}
			if got := li.SetLyricRenderer(tt.args.noteRenderer, tt.args.note); !assert.Equal(t, tt.want, got) {
				t.Errorf("lyricInteractor.SetLyricRenderer() = %v, want %v", got, tt.want)
			}

			assert.Equal(t, tt.wantWidth, tt.args.noteRenderer.Width, "note width after preparation")
			assert.Equal(t, tt.wantTakenFromLyric, tt.args.noteRenderer.IsLengthTakenFromLyric, "the width taken from lyric")
			assert.Equal(t, tt.wantLyric, tt.args.noteRenderer.Lyric, "the renderer lyric")
		})
	}
}

func Test_lyricInteractor_RenderLyrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv    func(*gomock.Controller) *canvas.MockCanvas
		measure []*entity.NoteRenderer
	}{
		{
			name: "empty measure",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='lyric'", "style='font-family:Caladea'")
				canv.EXPECT().Gend()
				return canv
			},
		},
		{
			name: "no lyric",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='lyric'", "style='font-family:Caladea'")
				canv.EXPECT().Gend()
				return canv
			},
			measure: []*entity.NoteRenderer{
				{
					Lyric: []entity.Lyric{},
				},
			},
		},
		{
			name: "no text",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='lyric'", "style='font-family:Caladea'")
				canv.EXPECT().Gend()
				return canv
			},
			measure: []*entity.NoteRenderer{
				{
					Lyric: []entity.Lyric{
						{
							Text:     []entity.Text{},
							Syllabic: musicxml.LyricSyllabicTypeBegin,
						},
					},
				},
			},
		},
		{
			name: "no prefix",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='lyric'", "style='font-family:Caladea'")
				canv.EXPECT().Text(50, 125, "U")
				canv.EXPECT().Text(60, 125, "nit")
				canv.EXPECT().Gend()
				return canv
			},
			measure: []*entity.NoteRenderer{
				{
					PositionX: 50,
					PositionY: 100,
					Lyric: []entity.Lyric{
						{
							Text: []entity.Text{
								{Value: "U"},
							},
							Syllabic: musicxml.LyricSyllabicTypeBegin,
						},
					},
				},
				{
					PositionX: 60,
					PositionY: 100,
					Lyric: []entity.Lyric{
						{
							Text: []entity.Text{
								{Value: "nit"},
							},
							Syllabic: musicxml.LyricSyllabicTypeEnd,
						},
					},
				},
			},
		},
		{
			name: "no prefix - and notes - with elsion",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='lyric'", "style='font-family:Caladea'")
				canv.EXPECT().Text(44, 125, "*Test")
				canv.EXPECT().Text(60, 125, "ing")
				canv.EXPECT().Qbez(59, 127, 66, 133, 73, 127, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Gend()
				return canv
			},
			measure: []*entity.NoteRenderer{
				{
					PositionX: 50,
					PositionY: 100,
					Lyric: []entity.Lyric{
						{
							Text: []entity.Text{
								{Value: "*T"},
								{Value: "es", Underline: 1},
								{Value: "t"},
							},
							Syllabic: musicxml.LyricSyllabicTypeBegin,
						},
					},
				},
				{
					PositionX: 60,
					PositionY: 100,
					Lyric: []entity.Lyric{
						{
							Text: []entity.Text{
								{Value: "ing"},
							},
							Syllabic: musicxml.LyricSyllabicTypeEnd,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var li lyricInteractor
			canv := canvas.Canvas(nil)
			if tt.canv != nil {
				canv = tt.canv(ctrl)
			}
			li.RenderLyrics(context.Background(), canv, tt.measure)
		})
	}
}
