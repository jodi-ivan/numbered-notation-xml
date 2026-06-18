package text

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
	"github.com/stretchr/testify/assert"
)

func Test_textInteractor_NoteHasText(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		note *entity.NoteRenderer
		t    []string
		want bool
	}{
		{
			name: "empty request and empty note",
			note: &entity.NoteRenderer{},
			want: false,
		},
		{
			name: "has request and empty note",
			note: &entity.NoteRenderer{},
			t:    []string{"one"},
			want: false,
		},
		{
			name: "has request and in the note",
			note: &entity.NoteRenderer{
				MeasureText: []musicxml.MeasureText{
					{
						Text: "one",
					},
				},
			},
			t:    []string{"one"},
			want: true,
		},
		{
			name: "all of the request satisfied",
			note: &entity.NoteRenderer{
				MeasureText: []musicxml.MeasureText{
					{Text: "one"},
					{Text: "two"},
					{Text: "three"},
				},
			},
			t:    []string{"one", "two"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ti textInteractor
			got := ti.NoteHasText(tt.note.MeasureText, tt.t...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_textInteractor_SetMeasureTextRenderer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		noteRenderer    *entity.NoteRenderer
		note            musicxml.Note
		getContext      func() context.Context
		isLastNote      bool
		want            bool
		wantMeasureText []musicxml.MeasureText
	}{
		{
			name:         "No text to set",
			note:         musicxml.Note{},
			noteRenderer: &entity.NoteRenderer{},
			want:         false,
		},
		// set text not on last and affected margin bottom
		{
			name: "not last note and affecting the margin bottom",
			note: musicxml.Note{
				MeasureText: []musicxml.MeasureText{
					{Text: "Refrein"},
				},
			},
			noteRenderer: &entity.NoteRenderer{},
			want:         true,
			wantMeasureText: []musicxml.MeasureText{
				{Text: "Refrein", TextAlignment: musicxml.TextAlignmentLeft},
			},
		},
		// set text on last and NOT affected margin bottom
		{
			name: "not last note and affecting the margin bottom",
			note: musicxml.Note{
				MeasureText: []musicxml.MeasureText{
					{Text: "unittest"},
				},
			},
			isLastNote:   true,
			noteRenderer: &entity.NoteRenderer{},
			want:         false,
			wantMeasureText: []musicxml.MeasureText{
				{Text: "unittest", TextAlignment: musicxml.TextAlignmentRight},
			},
		},
		// set text but removed due to param
		{
			name: "set text but removed due to param",
			note: musicxml.Note{
				MeasureText: []musicxml.MeasureText{
					{Text: "bait 2 disini"},
				},
			},
			noteRenderer:    &entity.NoteRenderer{},
			want:            false,
			wantMeasureText: []musicxml.MeasureText{},
			getContext: func() context.Context {
				param := &params.Param{
					Verse: 2,
				}

				return params.NewParamContext(context.Background(), param)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ti textInteractor
			ctx := context.Background()
			if tt.getContext != nil {
				ctx = tt.getContext()
			}
			got := ti.SetMeasureTextRenderer(ctx, tt.noteRenderer, tt.note, tt.isLastNote)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantMeasureText, tt.noteRenderer.MeasureText, "Measure text after assert")
		})
	}
}

func Test_textInteractor_RenderMeasureText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name      string // description of this test case
		y         int
		canv      func() *canvas.MockCanvasTestify
		lyric     func(c *gomock.Controller) *lyric.MockLyric
		notes     []*entity.NoteRenderer
		linestaff []lines.LineStaff
	}{
		{
			name: "no text",
			notes: []*entity.NoteRenderer{
				{},
				{},
			},
			canv: func() *canvas.MockCanvasTestify {
				return canvas.NewMockCanvasTestify(t)
			},
		},
		// with fermata  default text (refrein or fine)
		{
			name: "with fermata and default text (refrein or fine)",
			notes: []*entity.NoteRenderer{
				{
					Fermata: &musicxml.Femata{Type: musicxml.FermataTypeUpright},
					MeasureText: []musicxml.MeasureText{
						{Text: "Refrein"},
					},
					PositionX: 100,
				},
			},
			y: 100,
			canv: func() *canvas.MockCanvasTestify {
				res := canvas.NewMockCanvasTestify(t)
				res.EXPECT().Text(100, 60, "Refrein", []string{`style="font-style:italic"`})
				return res
			},
		},
		// with barline ending not default text and right alignment
		{
			name: "with barline ending not default text and right alignment",
			notes: []*entity.NoteRenderer{
				{
					Barline: &musicxml.Barline{
						Ending: &musicxml.BarlineEnding{},
					},
				},
				{
					MeasureText: []musicxml.MeasureText{
						{Text: "unittest", TextAlignment: musicxml.TextAlignmentRight},
					},
					PositionX: 100,
				},
			},
			y: 100,
			lyric: func(c *gomock.Controller) *lyric.MockLyric {
				res := lyric.NewMockLyric(c)
				res.EXPECT().CalculateLyricWidth("unittest").Return(16.0)
				return res
			},
			canv: func() *canvas.MockCanvasTestify {
				res := canvas.NewMockCanvasTestify(t)
				res.EXPECT().Text(734, 67, "unittest", []string{`style="font-style:italic;font-size:65%;font-weight:bold"`})
				return res
			},
		},
		// with linestaff  under placement? why

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ti textInteractor
			if tt.lyric != nil {
				ti.Lyric = tt.lyric(ctrl)
			}
			canv := tt.canv()
			ti.RenderMeasureText(context.Background(), tt.y, canv, tt.notes, tt.linestaff...)
			canv.AssertExpectations(t)
		})
	}
}
