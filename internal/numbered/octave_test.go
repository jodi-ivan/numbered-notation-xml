package numbered

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func Test_numberedInteractor_RenderOctave(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		note *entity.NoteRenderer
	}
	tests := []struct {
		name string
		args args

		initCanvas func(*gomock.Controller) *canvas.MockCanvas
	}{
		{
			name: "no octave",
			args: args{
				note: &entity.NoteRenderer{
					Octave: 0,
				},
			},
			initCanvas: func(c *gomock.Controller) *canvas.MockCanvas { return nil },
		},
		{
			name: "has octave up",
			args: args{
				note: &entity.NoteRenderer{
					PositionX: 50,
					PositionY: 150,
					Octave:    1,
				},
			},
			initCanvas: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Circle(55, 135, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
				return canv
			},
		},
		{
			name: "has octave down",
			args: args{
				note: &entity.NoteRenderer{
					PositionX: 50,
					PositionY: 150,
					Octave:    -1,
				},
			},
			initCanvas: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Circle(55, 155, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
				return canv
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ni := numberedInteractor{}
			pos := entity.Coordinate{X: float64(tt.args.note.PositionX), Y: float64(tt.args.note.PositionY)}
			ni.RenderOctave(context.Background(), tt.initCanvas(ctrl), tt.args.note.Octave, pos)
		})
	}
}
