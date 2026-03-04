package numbered

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func Test_numberedInteractor_RenderStrikethrough(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv          func(c *gomock.Controller) *canvas.MockCanvas
		strikethrough bool
		pos           entity.Coordinate
	}{
		{
			name: "nothing is happening",
		},
		{
			name:          "strikethrough",
			strikethrough: true,
			pos:           entity.Coordinate{X: 80, Y: 100},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Line(90, 84, 80, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.45")
				return canv
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ni numberedInteractor
			var canv *canvas.MockCanvas
			if tt.canv != nil {
				canv = tt.canv(ctrl)
			}
			ni.RenderStrikethrough(context.Background(), canv, tt.strikethrough, tt.pos)
		})
	}
}
