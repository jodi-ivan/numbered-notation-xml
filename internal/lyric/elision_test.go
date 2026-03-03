package lyric

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func Test_lyricInteractor_RenderElision(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		initCanvas func(*gomock.Controller) *canvas.MockCanvas

		text      []entity.Text
		lyricPart int
		pos       entity.Coordinate
	}{
		{
			name: "no text",
		},
		{
			name: "no underline",
			text: []entity.Text{
				entity.Text{Value: "unit"},
				entity.Text{Value: "test"},
			},
		},
		{
			name: "no underline",
			text: []entity.Text{
				entity.Text{Value: "unit"},
				entity.Text{Value: "test", Underline: 1},
			},
			initCanvas: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Qbez(107, 128, 119, 134, 131, 128, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				return canv
			},
			pos: entity.Coordinate{
				X: 80,
				Y: 100,
			},
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		t.Run(tt.name, func(t *testing.T) {
			var li lyricInteractor
			var canv *canvas.MockCanvas
			if tt.initCanvas != nil {
				canv = tt.initCanvas(ctrl)
			}
			li.RenderElision(context.Background(), canv, tt.text, tt.lyricPart, tt.pos)
		})
	}
}
