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
		notes []*entity.NoteRenderer
	}
	tests := []struct {
		name string
		args args

		initCanvas func(*gomock.Controller) *canvas.MockCanvas
	}{
		{
			name: "no octave",
			args: args{
				notes: []*entity.NoteRenderer{},
			},
			initCanvas: func(c *gomock.Controller) *canvas.MockCanvas { return nil },
		},
		{
			name: "has octave",
			args: args{
				notes: []*entity.NoteRenderer{
					&entity.NoteRenderer{
						PositionX: 50,
						PositionY: 150,
						Octave:    1,
					},
					&entity.NoteRenderer{
						PositionX: 60,
						PositionY: 150,
						Octave:    -1,
					},
				},
			},
			initCanvas: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='octaves'")
				canv.EXPECT().Circle(55, 135, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
				canv.EXPECT().Circle(65, 155, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.5")
				canv.EXPECT().Gend()
				return canv
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ni := numberedInteractor{}
			ni.RenderOctave(context.Background(), tt.initCanvas(ctrl), tt.args.notes)
		})
	}
}
