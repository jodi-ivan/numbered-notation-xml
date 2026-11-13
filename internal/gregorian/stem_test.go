package gregorian

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/stretchr/testify/assert"
)

func Test_renderStem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		canv      func(c *gomock.Controller) *canvas.MockCanvas
		lines     [5]int
		direction int
		start     []entity.CoordinateWithNoteLength
		end       []entity.CoordinateWithNoteLength
		want      StemInfo
	}{
		{
			name: "Simple Up Stem (Direction 1)",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				// Expected calculation: additional = -10.0. y3 = 10. y3+additional = 0.0
				res.EXPECT().LineFloat64(10.0, 20.0, 10.0, 0.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: 1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 20), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: -10.0,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(0, 4),
				HighestYPosition:   entity.NewCoordinate(0, 0),
			},
		},
		{
			name: "Simple Down Stem (Direction -1)",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				// Expected calculation: additional = 10.0. y3 = 10. y3+additional = 20.0
				res.EXPECT().LineFloat64(10.0, 20.0, 10.0, 20.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: -1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 20), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: 10.0,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(10, 20.0),
				HighestYPosition:   entity.NewCoordinate(0, 0.0),
			},
		},
		{
			name: "Stem with Horizontal Line Intersection and Note Length Adjustment",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				// Expected calculation: additional = 0.5. y3 = 2.5. y3+additional = 3.0
				res.EXPECT().LineFloat64(10.0, 20.0, 10.0, 3.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: 1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 20), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 2), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: 0.5,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(0, 4),
				HighestYPosition:   entity.NewCoordinate(0, 0.0),
			},
		},
		// Test Case: Max Rise Clamping (Up Stem)
		{
			name: "Max Rise Clamping Up Stem",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				res.EXPECT().LineFloat64(10.0, 10.0, 10.0, 30.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: 1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 30), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: 0.0,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(10, 30),
				HighestYPosition:   entity.NewCoordinate(0, 0),
			},
		},
		// Test Case: Max Rise Clamping (Down Stem)
		{
			name: "Max Rise Clamping Down Stem",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				// y1=0, y2=30. diffY=30. maxRise=16. y2 clamped to 0+16=16. clampY2 = 14.
				// y3 = 0. additional = 0.
				res.EXPECT().LineFloat64(10.0, 0.0, 10.0, 30.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: -1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 0), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 30), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: 0.0,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(10, 30),
				HighestYPosition:   entity.NewCoordinate(0, 0),
			},
		},
		// Test Case: Min Distance Clamping (Up Stem)
		{
			name: "Min Distance Clamping Up Stem",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				res.EXPECT().LineFloat64(10.0, 10.0, 10.0, -10.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: 1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: -20,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(00, 4.0),
				HighestYPosition:   entity.NewCoordinate(10, -10.0),
			},
		},
		// Test Case: Min Distance Clamping (Down Stem)
		{
			name: "Min Distance Clamping Down Stem",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				res.EXPECT().LineFloat64(10.0, 10.0, 10.0, 30.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: -1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: 20.0,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(10, 30),
				HighestYPosition:   entity.NewCoordinate(0, 0),
			},
		},
		// Test Case: Reset Condition (Up Stem)
		{
			name: "Reset Condition Triggered Up Stem",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				res.EXPECT().LineFloat64(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: 1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 30), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 5), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: 2.5,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(10, 7.5),
				HighestYPosition:   entity.NewCoordinate(0, 0),
			},
		},
		// Test Case: Horizontal Line Intersection
		{
			name: "Horizontal Line Intersection Adjustment",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				// y1=10, y2=10. lines contains 10. direction=1. additional += 2.5.
				// y3 = 10. y3+additional = 12.5.
				res.EXPECT().LineFloat64(10.0, 10.0, 10.0, -10.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: 1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: -20,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(0, 4),
				HighestYPosition:   entity.NewCoordinate(10, -10.0),
			},
		},
		// Test Case: Staff Line Intersection Adjustment (Quarter Note)
		{
			name: "Staff Line Intersection Adjustment Quarter Note",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				res.EXPECT().LineFloat64(10.0, 10.0, 10.0, -10.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: 1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: -20.0,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(0, 4),
				HighestYPosition:   entity.NewCoordinate(10, -10),
			},
		},
		// Test Case: Staff Line Intersection Adjustment (Half Note)
		{
			name: "Staff Line Intersection Adjustment Half Note",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)
				// y3 = 10. lines contains 10. NoteLengthHalf. direction=1. y3 += 0.5.
				// additional = 0. y3+additional = 10.5.
				res.EXPECT().LineFloat64(10.0, 10.0, 10.0, -10.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: 1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthHalf, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthHalf, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: -20.0,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(0, 4),
				HighestYPosition:   entity.NewCoordinate(10, -10.0),
			},
		},
		// Test Case: Multiple Notes (Interpolation)
		{
			name: "Multiple Notes Interpolation",
			canv: func(c *gomock.Controller) *canvas.MockCanvas {
				res := canvas.NewMockCanvas(c)

				res.EXPECT().LineFloat64(10.0, 20.0, 10.0, -10.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				res.EXPECT().LineFloat64(20.0, 10.0, 20.0, -10.0, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
				return res
			},
			lines:     [5]int{0, 1, 2, 3, 4},
			direction: 1,
			start: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 20), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
				{Coordinate: entity.NewCoordinate(20, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			end: []entity.CoordinateWithNoteLength{
				{Coordinate: entity.NewCoordinate(10, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
				{Coordinate: entity.NewCoordinate(20, 10), NoteLength: musicxml.NoteLengthQuarter, Beam: nil},
			},
			want: StemInfo{
				LengthCompensation: -20.0,
				ClampY1:            0,
				ClampY2:            0,
				LowestYPosition:    entity.NewCoordinate(0, 4),
				HighestYPosition:   entity.NewCoordinate(10, -10.0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderStem(tt.canv(ctrl), tt.lines, tt.direction, tt.start, tt.end)
			assert.Equal(t, got, tt.want)
		})
	}
}
