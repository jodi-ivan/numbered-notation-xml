package rhythm

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
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
				canv.EXPECT().Qbez(55, 105, 80, 118, 105, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
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
				canv.EXPECT().Qbez(57, 107, 130, 123, 203, 107, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Gend()
				return canv
			},
		},
		{
			name: "with octave dashed",
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
						LineType: musicxml.NoteSlurLineTypeDashed,
					},
				},
			},
			initCanvasMock: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='slurties'")
				canv.EXPECT().Qbez(57, 107, 130, 123, 203, 107, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1;stroke-dasharray:3.5 3.8;stroke-dashoffset:1.750000;")
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

func Test_rhythmInteractor_RenderSlurTies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		// Named input parameters for target function.
		canv         func(*gomock.Controller) *canvas.MockCanvas
		notes        []*entity.NoteRenderer
		maxXPosition float64
	}{
		{
			name: "no ties and no slur",
			notes: []*entity.NoteRenderer{
				&entity.NoteRenderer{},
			},
		},
		{
			name: "ties only",
			notes: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					PositionX: 80,
					PositionY: 100,
					Tie: &entity.Slur{
						Number: 1,
						Type:   musicxml.NoteSlurTypeStart,
					},
				},
				&entity.NoteRenderer{
					PositionX: 150,
					PositionY: 100,
					Tie: &entity.Slur{
						Number: 1,
						Type:   musicxml.NoteSlurTypeStop,
					},
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='slurties'")
				canv.EXPECT().Qbez(85, 105, 120, 118, 155, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Gend()
				return canv
			},
		},
		{
			name: "Slur only",
			notes: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					PositionX: 80,
					PositionY: 100,
					Slur: map[int]entity.Slur{
						1: entity.Slur{
							Number: 1,
							Type:   musicxml.NoteSlurTypeStart,
						},
					},
				},
				&entity.NoteRenderer{
					PositionX: 150,
					PositionY: 100,
					Slur: map[int]entity.Slur{
						1: entity.Slur{
							Number: 1,
							Type:   musicxml.NoteSlurTypeStop,
						},
					},
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='slurties'")
				canv.EXPECT().Qbez(87, 105, 120, 118, 153, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Gend()
				return canv
			},
		},
		{
			name: "Slur only no start",
			notes: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					PositionX: 150,
					PositionY: 100,
					Slur: map[int]entity.Slur{
						1: entity.Slur{
							Number: 1,
							Type:   musicxml.NoteSlurTypeStop,
						},
					},
				},
			},
			maxXPosition: 200,
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='slurties'")
				canv.EXPECT().Qbez(135, 105, 144, 111, 153, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Gend()
				return canv
			},
		},
		{
			name: "Slur only no end",
			notes: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					PositionX: 100,
					PositionY: 100,
					Slur: map[int]entity.Slur{
						1: entity.Slur{
							Number: 1,
							Type:   musicxml.NoteSlurTypeStart,
						},
					},
				},
			},
			maxXPosition: 200,
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='slurties'")
				canv.EXPECT().Qbez(107, 105, 154, 118, 200, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Gend()
				return canv
			},
		},
		{
			name: "Slur with hops",
			notes: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					PositionX: 80,
					PositionY: 100,
					Slur: map[int]entity.Slur{
						1: entity.Slur{
							Number: 1,
							Type:   musicxml.NoteSlurTypeStart,
						},
					},
				},
				&entity.NoteRenderer{
					PositionX: 110,
					PositionY: 100,
					Slur: map[int]entity.Slur{
						1: entity.Slur{
							Number: 1,
							Type:   musicxml.NoteSlurTypeHop,
						},
					},
				},
				&entity.NoteRenderer{
					PositionX: 150,
					PositionY: 100,
					Slur: map[int]entity.Slur{
						1: entity.Slur{
							Number: 1,
							Type:   musicxml.NoteSlurTypeStop,
						},
					},
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='slurties'")
				canv.EXPECT().Qbez(87, 105, 100, 111, 113, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Qbez(117, 105, 135, 111, 153, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Gend()
				return canv
			},
		},
		{
			name: "Both slur and ties",
			notes: []*entity.NoteRenderer{
				&entity.NoteRenderer{
					PositionX: 80,
					PositionY: 100,
					Slur: map[int]entity.Slur{
						1: entity.Slur{
							Number: 1,
							Type:   musicxml.NoteSlurTypeStart,
						},
					},
				},
				&entity.NoteRenderer{
					PositionX: 150,
					PositionY: 100,
					Slur: map[int]entity.Slur{
						1: entity.Slur{
							Number: 1,
							Type:   musicxml.NoteSlurTypeStop,
						},
					},
				},
				&entity.NoteRenderer{
					PositionX: 160,
					PositionY: 100,
					Tie: &entity.Slur{
						Number: 1,
						Type:   musicxml.NoteSlurTypeStart,
					},
				},
				&entity.NoteRenderer{
					PositionX: 200,
					PositionY: 100,
					Tie: &entity.Slur{
						Number: 1,
						Type:   musicxml.NoteSlurTypeStop,
					},
				},
			},
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				canv := canvas.NewMockCanvas(c)
				canv.EXPECT().Group("class='slurties'").Times(2)
				canv.EXPECT().Qbez(165, 105, 185, 118, 205, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Qbez(87, 105, 120, 118, 153, 105, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1")
				canv.EXPECT().Gend().Times(2)
				return canv
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var canv *canvas.MockCanvas
			if tt.canv != nil {
				canv = tt.canv(ctrl)
			}
			var ri rhythmInteractor
			ri.RenderSlurTies(context.Background(), canv, tt.notes, tt.maxXPosition)
		})
	}
}
