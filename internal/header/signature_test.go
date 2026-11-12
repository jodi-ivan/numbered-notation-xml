package header

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	ubermock "go.uber.org/mock/gomock"
)

func Test_headerInteractor_RenderKeyandTimeSignatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uberCtrl := ubermock.NewController(t)
	defer uberCtrl.Finish()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv          func(*gomock.Controller) *canvas.MockCanvas
		lyricMock     func(*ubermock.Controller) *lyric.MockLyric
		key           keysig.KeySignature
		timeSignature timesig.TimeSignature
	}{
		{
			name: "BAU",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Text(50, 60, "do = c")
				canv.EXPECT().Text(175, 60, "4 ketuk")
				return canv
			},
			lyricMock: func(c *ubermock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth("do = c").Return(80.0)
				return li
			},
			key: keysig.NewKeySignature(context.Background(), []musicxml.Measure{
				{
					Attribute: &musicxml.Attribute{
						Key: &musicxml.KeySignature{
							Fifth: 0,
							Mode:  "Major",
						},
					},
				},
			}),
			timeSignature: timesig.NewTimeSignatures(context.Background(), []musicxml.Measure{
				{
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int "xml:\"beats\""
							BeatType int "xml:\"beat-type\""
						}{
							Beats:    4,
							BeatType: 4,
						},
					},
				},
			}),
		},
		{
			name: "MultipleKeySignatures",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				// Expect calls for multiple key signatures
				canv.EXPECT().Text(50, 60, gomock.Any())
				canv.EXPECT().Text(175, 60, "4 ketuk")
				return canv
			},
			lyricMock: func(c *ubermock.Controller) *lyric.MockLyric {
				li := lyric.NewMockLyric(c)
				li.EXPECT().CalculateLyricWidth(gomock.Any()).Return(80.0)
				return li
			},
			key: keysig.NewKeySignature(context.Background(), []musicxml.Measure{
				{
					Attribute: &musicxml.Attribute{
						Key: &musicxml.KeySignature{
							Fifth: 0,
							Mode:  "Major",
						},
					},
				},
				{
					Attribute: &musicxml.Attribute{
						Key: &musicxml.KeySignature{
							Fifth: 2,
							Mode:  "Major",
						},
					},
				},
				{
					Attribute: &musicxml.Attribute{
						Key: &musicxml.KeySignature{
							Fifth: 3,
							Mode:  "Major",
						},
					},
				},
			}),
			timeSignature: timesig.NewTimeSignatures(context.Background(), []musicxml.Measure{
				{
					Attribute: &musicxml.Attribute{
						Time: &struct {
							Beats    int "xml:\"beats\""
							BeatType int "xml:\"beat-type\""
						}{
							Beats:    4,
							BeatType: 4,
						},
					},
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var hi headerInteractor
			if tt.lyricMock != nil {
				hi.Lyric = tt.lyricMock(uberCtrl)
			}
			canv := canvas.Canvas(nil)
			if tt.canv != nil {
				canv = tt.canv(ctrl)
			}
			hi.RenderKeyandTimeSignatures(context.Background(), canv, tt.key, tt.timeSignature)
		})
	}
}
