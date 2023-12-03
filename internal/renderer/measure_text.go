package renderer

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderMeasureTopping(ctx context.Context, canv canvas.Canvas, notes []*entity.NoteRenderer) {
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
		musicxml.BarLineStyleRegular:    7,
		musicxml.BarLineStyleLightHeavy: 2,
	}

	offsetEnd := map[musicxml.BarLineStyle]int{
		musicxml.BarLineStyleLightHeavy: -3,
	}
	for _, note := range notes {
		if note.Barline != nil && note.Barline.Ending != nil {
			switch note.Barline.Ending.Type {
			case musicxml.BarlineEndingTypeStart:
				pairs = append(pairs, [2]entity.Coordinate{
					entity.Coordinate{
						X: float64(note.PositionX),
						Y: float64(note.PositionY),
					},
				})
				pairsData = append(pairsData, note.Barline.Ending.Number)
				beginTarget, ok := rightMeasureMap[note.MeasureNumber-1]
				if !ok {
					beginTarget = musicxml.BarLineStyleRegular
				}
				pairsBar = append(pairsBar, [2]musicxml.BarLineStyle{
					beginTarget,
				})
			case musicxml.BarlineEndingTypeStop, musicxml.BarlineEndingTypeDiscontinue:
				curr := pairs[len(pairs)-1]
				curr[1] = entity.Coordinate{
					X: float64(note.PositionX),
					Y: float64(note.PositionY),
				}

				pairs[len(pairs)-1] = curr

				pairsBar[len(pairsBar)-1][1] = note.Barline.BarStyle

			}
		}
	}
	if len(pairs) > 0 {
		canv.Group("class='staff-topping'")

		for i, pair := range pairs {
			x1 := int(math.Round(pair[0].X)) - offsetStart[pairsBar[i][0]]
			x2 := int(math.Round(pair[1].X)) - offsetEnd[pairsBar[i][1]]
			canv.Text(x1+3, int(math.Round(pair[0].Y))-12, pairsData[i], `style="font-weight:bold;font-size:90%"`)
			// vertical line at start
			canv.Line(
				x1,
				int(math.Round(pair[0].Y))-18,
				x1,
				int(math.Round(pair[1].Y))-25,
				"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
			)
			canv.Line(
				x1,
				int(math.Round(pair[0].Y))-25,
				x2,
				int(math.Round(pair[1].Y))-25,
				"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
			)
			// vertical line at end
			canv.Line(
				x2,
				int(math.Round(pair[0].Y))-18,
				x2,
				int(math.Round(pair[1].Y))-25,
				"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
			)
		}
		canv.Gend()
	}

}

func RenderMeasureText(ctx context.Context, canv canvas.Canvas, notes []*entity.NoteRenderer) {
	canv.Group("class='staff-text'")
	for _, note := range notes {
		if len(note.MeasureText) > 0 {
			sort.Slice(note.MeasureText, func(i, j int) bool {
				return note.MeasureText[i].RelativeY < note.MeasureText[j].RelativeY
			})

			for i, t := range note.MeasureText {
				xPos := note.PositionX
				if t.TextAlignment == musicxml.TextAlignmentRight {
					textLength := lyric.CalculateLyricWidth(t.Text)
					xPos = xPos - int(textLength)
				}
				canv.Text(xPos, note.PositionY-28-(i*-15), t.Text, `font-style="italic"`)
			}
		}

	}
	canv.Gend()

}

func RenderTuplet(ctx context.Context, canv canvas.Canvas, notes []*entity.NoteRenderer) {
	pairs := [][2]entity.Coordinate{}
	pairData := []int{}
	for _, n := range notes {
		if n.Tuplet != nil {

			switch n.Tuplet.Type {
			case musicxml.TupletTypeStart:
				pairs = append(pairs, [2]entity.Coordinate{
					entity.Coordinate{
						X: float64(n.PositionX),
						Y: float64(n.PositionY),
					},
				})
				pairData = append(pairData, n.TimeMofication.ActualNotes.Value)
			case musicxml.TupletTypeStop:
				curr := pairs[len(pairs)-1]
				curr[1] = entity.Coordinate{
					X: float64(n.PositionX),
					Y: float64(n.PositionY),
				}

				pairs[len(pairs)-1] = curr
			}
		}
	}

	if len(pairs) > 0 {

		canv.Group("class='tuplet'", `style="font-size:80%"`)
		for i, pair := range pairs {
			end := pair[1]
			start := pair[0]

			x := end.X - start.X

			canv.Text(int((start.X + (x / 2))), int(start.Y)-20, fmt.Sprintf("%d", pairData[i]))
		}
		canv.Gend()
	}
}
