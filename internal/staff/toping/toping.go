package toping

import (
	"context"
	"fmt"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/text"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type Toping interface {
	RenderRepeatMeasure(ctx context.Context, y int, canv canvas.Canvas, notes []*entity.NoteRenderer, std ...bool)
	RenderTuplet(ctx context.Context, y int, canv canvas.Canvas, notes []*entity.NoteRenderer)
	SetStaffLineDashRenderer(note *entity.NoteRenderer, dashes map[int]musicxml.DirectionDashesType)
	RenderStaffLineDash(notes []*entity.NoteRenderer, canv canvas.Canvas, y int, linestaff ...lines.LineStaff)
}

func NewToping() Toping {
	return &topingInteractor{}
}

type topingInteractor struct {
}

func (ti *topingInteractor) RenderTuplet(ctx context.Context, y int, canv canvas.Canvas, notes []*entity.NoteRenderer) {
	pairs := [][2]CoordinateWithTuplet{}
	pairData := []int{}
	for _, n := range notes {
		if n.Tuplet == nil {
			continue
		}

		switch n.Tuplet.Type {
		case musicxml.TupletTypeStart:
			pairs = append(pairs, [2]CoordinateWithTuplet{
				{
					Coordinate: entity.NewCoordinate(float64(n.PositionX), float64(y)),
					Tuplet:     *n.Tuplet,
				},
			})
			pairData = append(pairData, n.TimeModifications.ActualNotes.Value)
		case musicxml.TupletTypeStop:
			curr := pairs[len(pairs)-1]
			curr[1] = CoordinateWithTuplet{
				Coordinate: entity.NewCoordinate(float64(n.PositionX), float64(y)),
				Tuplet:     *n.Tuplet,
			}

			pairs[len(pairs)-1] = curr
		}
	}

	if len(pairs) > 0 {

		canv.Group("class='tuplet'", `style="font-size:80%"`)
		for i, pair := range pairs {
			end := pair[1]
			start := pair[0]

			x := start.X + ((end.X - start.X) / 2)
			if start.Tuplet.Braket == musicxml.BoolYes {
				canv.Qbez(
					int(start.X), int(end.Y)-22,
					int(x)+4, int(start.Y)-38,
					int(end.X)+8, int(end.Y)-22,
					"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.8",
				)
				canv.CenterRect(int(x)+4, int(start.Y)-26, 10, 12, "fill:white;stroke:none;")
			}
			canv.Text(int(x), int(start.Y)-22, fmt.Sprintf("%d", pairData[i]), "font-style:italic")

		}
		canv.Gend()
	}
}

func (ti *topingInteractor) RenderRepeatMeasure(ctx context.Context, y int, canv canvas.Canvas, notes []*entity.NoteRenderer, std ...bool) {
	pairs := [][2]entity.Coordinate{}
	pairsData := []string{}
	pairsBar := [][2]musicxml.BarLineStyle{}

	rightMeasureMap := map[int]musicxml.BarLineStyle{}

	for _, n := range notes {
		if n.Barline != nil && n.Barline.Location == musicxml.BarlineLocationRight {
			rightMeasureMap[n.MeasureNumber] = n.Barline.BarStyle
		}
	}

	offsetStart := map[musicxml.BarLineStyle]int{
		musicxml.BarLineStyleRegular:    2,
		musicxml.BarLineStyleLightHeavy: int(barline.GetBarlineWidth(musicxml.BarLineStyleLightHeavy)),
	}

	offsetEnd := map[musicxml.BarLineStyle]int{
		musicxml.BarLineStyleLightHeavy: -3,
	}
	for notePos, note := range notes {
		if note.Barline != nil && note.Barline.Ending != nil {
			switch note.Barline.Ending.Type {
			case musicxml.BarlineEndingTypeStart:
				pos := entity.NewCoordinate(float64(note.PositionX), float64(y))
				if notePos > 0 && notes[notePos-1].Barline != nil &&
					notes[notePos-1].Barline.BarStyle != musicxml.BarLineStyleNone {
					pos = entity.NewCoordinate(float64(notes[notePos-1].PositionX), float64(y))

				}
				pairs = append(pairs, [2]entity.Coordinate{pos})
				pairsData = append(pairsData, note.Barline.Ending.Number)
				beginTarget, ok := rightMeasureMap[note.MeasureNumber-1]
				if !ok {
					beginTarget = musicxml.BarLineStyleRegular
				}
				pairsBar = append(pairsBar, [2]musicxml.BarLineStyle{beginTarget})
			case musicxml.BarlineEndingTypeStop, musicxml.BarlineEndingTypeDiscontinue:
				curr := pairs[len(pairs)-1]
				curr[1] = entity.NewCoordinate(float64(note.PositionX), float64(y))

				pairs[len(pairs)-1] = curr

				pairsBar[len(pairsBar)-1][1] = note.Barline.BarStyle

			}
		}
	}

	y1Offset := 25
	if len(std) > 0 {
		y1Offset = 15
	}
	if len(pairs) > 0 {
		canv.Group("class='staff-topping'")

		for i, pair := range pairs {
			x1 := int(math.Round(pair[0].X)) + offsetStart[pairsBar[i][0]]
			x2 := int(math.Round(pair[1].X)) - offsetEnd[pairsBar[i][1]] - 5
			canv.Text(x1+3, int(math.Round(pair[0].Y))-22, pairsData[i], `style="font-weight:bold;font-size:60%;font-family:Old Standard TT;"`)
			// vertical line at start
			canv.Line(
				x1, int(math.Round(pair[0].Y))-y1Offset,
				x1, int(math.Round(pair[1].Y))-30,
				"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
			)
			canv.Line(
				x1, int(math.Round(pair[0].Y))-30,
				x2, int(math.Round(pair[1].Y))-30,
				"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
			)
			if i%2 == 1 {
				continue
			}
			// vertical line at end
			canv.Line(
				x2, int(math.Round(pair[0].Y))-y1Offset,
				x2, int(math.Round(pair[1].Y))-30,
				"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
			)
		}
		canv.Gend()
	}

}

func (ti *topingInteractor) SetStaffLineDashRenderer(note *entity.NoteRenderer, dashes map[int]musicxml.DirectionDashesType) {
	if dashes == nil {
		return
	}
	if note.MeasureDash == nil {
		note.MeasureDash = map[int]musicxml.DirectionDashesType{}
	}

	for num, dashType := range dashes {
		note.MeasureDash[num] = dashType
	}
}

func (ti *topingInteractor) RenderStaffLineDash(notes []*entity.NoteRenderer, canv canvas.Canvas, y int, linestaff ...lines.LineStaff) {
	dashSet := map[int][2]entity.Coordinate{}

	for notePos, note := range notes {
		noteMarginTop := 0.0

		if len(linestaff) > 0 {
			noteMarginTop = text.GetTextMarginBottom(linestaff[0], notes, notePos)
		}

		for num, dashType := range note.MeasureDash {

			pair, ok := dashSet[num]
			if !ok {
				pair = [2]entity.Coordinate{}
			}

			loc := entity.NewCoordinate(float64(note.PositionX), float64(y))

			if len(linestaff) > 0 {
				loc.Y -= noteMarginTop
			}

			if dashType == musicxml.DirectionDashesTypeStop {
				loc.X = float64(notes[notePos-1].PositionX) + constant.LOWERCASE_LENGTH
			}

			switch dashType {
			case musicxml.DirectionDashesTypeStart:
				pair[0] = loc
			case musicxml.DirectionDashesTypeStop:
				loc.Y = pair[0].Y
				pair[1] = loc

			}
			dashSet[num] = pair
		}
	}

	for _, pair := range dashSet {
		if pair[1].X == 0 {
			n := notes[len(notes)-1]
			if n.Barline != nil {
				n = notes[len(notes)-2]
			}
			pair[1].X, pair[1].Y = float64(n.PositionX)+4, pair[0].Y
		}
		canv.Line(int(pair[0].X)+constant.LOWERCASE_LENGTH, int(pair[0].Y)-25, int(pair[1].X), int(pair[1].Y)-25, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1;stroke-dasharray:4 8;")
	}

}
