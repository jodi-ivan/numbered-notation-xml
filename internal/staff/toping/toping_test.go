package toping

import (
	"context"
	"testing"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func Test_topingInteractor_RenderTuplet(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		y     int
		canv  func() *canvas.MockCanvasTestify
		notes []*entity.NoteRenderer
	}{
		// no tuplet
		{
			name: "no tuplet",
			y:    100,
			canv: func() *canvas.MockCanvasTestify {
				return canvas.NewMockCanvasTestify(t)
			},
			notes: []*entity.NoteRenderer{
				{}, {},
			},
		},
		// with tuplet without bracket
		{
			name: "with tuplet without bracket",
			y:    100,
			notes: []*entity.NoteRenderer{
				{
					PositionX: 100,
					Tuplet:    &musicxml.Tuplet{Type: musicxml.TupletTypeStart},
					TimeModifications: &musicxml.TimeModification{
						ActualNotes: musicxml.ChardataInt{Value: 1},
					},
				},
				{
					PositionX: 120,
					Tuplet:    &musicxml.Tuplet{Type: musicxml.TupletTypeStop},
					TimeModifications: &musicxml.TimeModification{
						ActualNotes: musicxml.ChardataInt{Value: 2},
					},
				},
			},
			canv: func() *canvas.MockCanvasTestify {
				res := canvas.NewMockCanvasTestify(t)
				res.EXPECT().Group([]string{"class='tuplet'", `style="font-size:80%"`})
				res.EXPECT().Text(110, 78, "1", []string{"font-style:italic"})
				res.EXPECT().Gend()
				return res
			},
		},
		// with tuplet with bracket
		{
			name: "with tuplet with bracket",
			y:    100,
			notes: []*entity.NoteRenderer{
				{
					PositionX: 100,
					Tuplet:    &musicxml.Tuplet{Type: musicxml.TupletTypeStart, Braket: musicxml.BoolYes},
					TimeModifications: &musicxml.TimeModification{
						ActualNotes: musicxml.ChardataInt{Value: 3},
					},
				},
				{
					PositionX: 120,
					Tuplet:    &musicxml.Tuplet{Type: musicxml.TupletTypeStop},
				},
			},
			canv: func() *canvas.MockCanvasTestify {
				res := canvas.NewMockCanvasTestify(t)
				res.EXPECT().Group([]string{"class='tuplet'", `style="font-size:80%"`})
				res.EXPECT().Text(110, 78, "3", []string{"font-style:italic"})
				res.EXPECT().Qbez(100, 78, 114, 62, 128, 78, []string{"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.8"})
				res.EXPECT().CenterRect(114, 74, 10, 12, []string{"fill:white;stroke:none;"})
				res.EXPECT().Gend()
				return res
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ti topingInteractor
			canv := tt.canv()
			ti.RenderTuplet(context.Background(), tt.y, canv, tt.notes)
			canv.AssertExpectations(t)
		})
	}
}

func Test_topingInteractor_RenderRepeatMeasure(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		y     int
		canv  func() *canvas.MockCanvasTestify
		notes []*entity.NoteRenderer
		std   []bool
	}{
		// no repeat
		{
			name: "no repeat",
			y:    100,
			canv: func() *canvas.MockCanvasTestify {
				return canvas.NewMockCanvasTestify(t)
			},
			notes: []*entity.NoteRenderer{
				{}, {Barline: &musicxml.Barline{BarStyle: musicxml.BarLineStyleRegular}},
			},
		},
		// has repeat
		{
			name: "has repeat",
			y:    100,
			canv: func() *canvas.MockCanvasTestify {
				res := canvas.NewMockCanvasTestify(t)
				res.EXPECT().Group([]string{"class='staff-topping'"})
				res.EXPECT().Text(105, 78, "1", []string{`style="font-weight:bold;font-size:60%;font-family:Old Standard TT;"`})
				res.EXPECT().Line(102, 75, 102, 70, []string{"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1"})
				res.EXPECT().Line(102, 70, 105, 70, []string{"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1"})
				res.EXPECT().Line(105, 75, 105, 70, []string{"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1"})

				res.EXPECT().Gend()
				return res
			},
			notes: []*entity.NoteRenderer{
				{
					PositionX: 100,
					Barline: &musicxml.Barline{
						BarStyle: musicxml.BarLineStyleRegular,
						Ending: &musicxml.BarlineEnding{
							Number: "1",
							Type:   musicxml.BarlineEndingTypeStart,
						},
					},
				},
				{
					PositionX: 110,
					Barline: &musicxml.Barline{
						BarStyle: musicxml.BarLineStyleRegular,
						Ending: &musicxml.BarlineEnding{
							Number: "2",
							Type:   musicxml.BarlineEndingTypeStop,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ti topingInteractor
			canv := tt.canv()
			ti.RenderRepeatMeasure(context.Background(), tt.y, canv, tt.notes, tt.std...)
			canv.AssertExpectations(t)
		})
	}
}

func Test_topingInteractor_RenderStaffLineDash(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		notes     []*entity.NoteRenderer
		canv      func() *canvas.MockCanvasTestify
		y         int
		linestaff []lines.LineStaff
	}{
		// no line
		{
			name: "no line",
			notes: []*entity.NoteRenderer{
				{}, {}, {},
			},
			canv: func() *canvas.MockCanvasTestify {
				return canvas.NewMockCanvasTestify(t)
			},
		},
		{
			name: "with line",
			notes: []*entity.NoteRenderer{
				{
					PositionX: 100,
					MeasureDash: map[int]musicxml.DirectionDashesType{
						1: musicxml.DirectionDashesTypeStart,
					},
				}, {},
				{
					PositionX: 110,
					MeasureDash: map[int]musicxml.DirectionDashesType{
						1: musicxml.DirectionDashesTypeStop,
					},
					Barline: &musicxml.Barline{},
				},
			},
			y: 100,
			canv: func() *canvas.MockCanvasTestify {
				res := canvas.NewMockCanvasTestify(t)

				res.EXPECT().Line(115, 75, 15, 75, []string{"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1;stroke-dasharray:4 8;"})

				return res
			},
		},
		{
			name: "with line with no end",
			notes: []*entity.NoteRenderer{
				{
					PositionX: 100,
					MeasureDash: map[int]musicxml.DirectionDashesType{
						1: musicxml.DirectionDashesTypeStart,
					},
				}, {PositionX: 500}, {Barline: &musicxml.Barline{}},
			},
			y: 100,
			canv: func() *canvas.MockCanvasTestify {
				res := canvas.NewMockCanvasTestify(t)

				res.EXPECT().Line(115, 75, 504, 75, []string{"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1;stroke-dasharray:4 8;"})

				return res
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ti topingInteractor
			canv := tt.canv()
			ti.RenderStaffLineDash(tt.notes, canv, tt.y, tt.linestaff...)
			canv.AssertExpectations(t)
		})
	}
}
