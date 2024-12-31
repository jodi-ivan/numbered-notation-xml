package rhythm

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func Test_rhythmInteractor_RenderBezier(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		set []SlurBezier
	}
	tests := []struct {
		name           string
		args           args
		initCanvasMock func(*gomock.Controller) *canvas.MockCanvas
	}{
		{
			name: "no set",
			args: args{},
			initCanvasMock: func(c *gomock.Controller) *canvas.MockCanvas {
				return nil
			},
		},
		{
			name: "with no octave",
			args: args{
				set: []SlurBezier{
					SlurBezier{
						Start: CoordinateWithOctave{
							Coordinate: entity.Coordinate{
								X: 50,
								Y: 100,
							},
						},
						End: CoordinateWithOctave{
							Coordinate: entity.Coordinate{
								X: 100,
								Y: 100,
							},
						},
					},
				},
			},
			initCanvasMock: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='slurties'")
				canv.EXPECT().Qbez(55, 105, 80, 113, 105, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.5")
				canv.EXPECT().Gend()
				return canv
			},
		},
		{
			name: "with octave",
			args: args{
				set: []SlurBezier{
					SlurBezier{
						Start: CoordinateWithOctave{
							Coordinate: entity.Coordinate{
								X: 50,
								Y: 100,
							},
							Octave: -1,
						},
						End: CoordinateWithOctave{
							Coordinate: entity.Coordinate{
								X: 200,
								Y: 100,
							},
							Octave: -1,
						},
					},
				},
			},
			initCanvasMock: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='slurties'")
				canv.EXPECT().Qbez(58, 108, 130, 118, 202, 108, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.5")
				canv.EXPECT().Gend()
				return canv
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ri := &rhythmInteractor{}
			ri.RenderBezier(tt.args.set, tt.initCanvasMock(ctrl))
		})
	}
}
