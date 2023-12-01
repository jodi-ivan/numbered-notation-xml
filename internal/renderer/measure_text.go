package renderer

import (
	"context"
	"fmt"
	"sort"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

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
				canv.Text(xPos, note.PositionY-25-(i*-15), t.Text, `font-style="italic"`)
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
