package renderer

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderVerse(ctx context.Context, canv canvas.Canvas, y int, verses []repository.HymnVerse) VerseInfo {
	canv.Group("class='verse'", "style='font-family:Caladea'")
	initialY := y

	combine := [][2]entity.Coordinate{}

	for _, verse := range verses {

		whole := [][]lyric.LyricWordVerse{}

		json.Unmarshal([]byte(verse.Content.String), &whole)

		lineLength := float64(0)

		blob := []string{}

		for _, line := range whole {
			lineText := ""
			for iLine, word := range line {
				lineText = lineText + " " + word.Word
				wordPart := ""
				for _, p := range word.Breakdown {

					if p.Combine {

						startPosition := lyric.CalculateLyricWidth(lineText) - lyric.CalculateLyricWidth(word.Word) + lyric.CalculateLyricWidth(wordPart)
						combine = append(combine, [2]entity.Coordinate{
							entity.Coordinate{
								X: startPosition,
								Y: float64(initialY + (iLine * 25)),
							},
							entity.Coordinate{
								X: startPosition + lyric.CalculateLyricWidth(p.Text[1:]),
								Y: float64(initialY + (iLine * 25)),
							},
						})
					}
					wordPart = wordPart + p.Text
				}
			}
			lineLength = math.Max(lineLength, lyric.CalculateLyricWidth(lineText)+float64(4*len(line)))
			blob = append(blob, lineText)
		}

		x := int(math.Round((constant.LAYOUT_WIDTH / 2) - (lineLength / 2)))
		canv.Group()
		canv.Text(x-int(lyric.CalculateLyricWidth(fmt.Sprintf("%d. ", verse.VerseNum.Int32))), initialY, fmt.Sprintf("%d. ", verse.VerseNum.Int32))
		for _, l := range blob {
			canv.Text(x, y, l)
			y += 25
		}
		canv.Gend()

		canv.Group()

		for _, c := range combine {
			canv.Qbez(
				int(c[0].X)+x, int(c[0].Y+2),
				int(c[0].X)+x+(int(c[1].X-c[0].X)/2), int(c[0].Y+7),
				int(c[1].X)+x, int(c[1].Y+2),
				"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
			)
		}
		canv.Gend()
	}

	canv.Gend()
	return VerseInfo{
		MarginBottom: y,
	}
}
