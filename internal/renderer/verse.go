package renderer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderVerse(ctx context.Context, canv canvas.Canvas, y int, verses []repository.HymnVerse) VerseInfo {
	canv.Group("class='verses'", "style='font-family:Caladea'")

	allVerse := map[int][]string{}
	lineLength := float64(0)
	allCombine := map[int][][2]entity.Coordinate{}

	for _, verse := range verses {

		combine := [][2]entity.Coordinate{}
		whole := [][]lyric.LyricWordVerse{}

		err := json.Unmarshal([]byte(verse.Content.String), &whole)
		if err != nil {
			log.Println("[RenderVerse] failed to unmarshal, err ", err)
		}

		blob := []string{}
		for iLine, line := range whole {
			lineText := ""
			for _, word := range line {
				wordPart := ""
				for _, p := range word.Breakdown {

					if p.Combine {
						x1 := float64(0)
						x2 := float64(0)
						for _, v := range p.Breakdown {
							if v.Underline {
								x2 = x1 + lyric.CalculateLyricWidth(v.Text)
								break
							} else {
								x1 += lyric.CalculateLyricWidth(v.Text)

							}
						}

						startPosition := lyric.CalculateLyricWidth(lineText) + lyric.CalculateLyricWidth(wordPart)
						combine = append(combine, [2]entity.Coordinate{
							entity.Coordinate{
								X: startPosition + x1,
								Y: float64(iLine),
							},
							entity.Coordinate{
								X: startPosition + x2,
								Y: float64(iLine),
							},
						})
					}
					wordPart = wordPart + p.Text
				}
				lineText = lineText + " " + word.Word

			}
			lineLength = math.Max(lineLength, lyric.CalculateLyricWidth(lineText))
			blob = append(blob, lineText)
		}
		allVerse[int(verse.VerseNum.Int32)] = blob
		allCombine[int(verse.VerseNum.Int32)] = combine
	}

	x := int(math.Round((constant.LAYOUT_WIDTH / 2) - (lineLength / 2)))
	totalVerse := len(allVerse)

	for i := 1; i < totalVerse+1; i++ {

		canv.Group("class='verse'")
		yVerse := y

		currentVerse := allVerse[i+1]

		// number verse
		canv.Text(x-5-int(lyric.CalculateLyricWidth(fmt.Sprintf("%d. ", i+1))), y, fmt.Sprintf("%d. ", i+1))
		for _, liveVerse := range currentVerse {
			canv.Text(x, y, liveVerse)
			y += 25

		}
		canv.Group()

		for _, c := range allCombine[i+1] {
			canv.Qbez(
				int(c[0].X)+x, int(c[0].Y*25)+2+yVerse,
				int(c[0].X)+x+(int(c[1].X-c[0].X)/2), yVerse+7+(int(c[0].Y)*25),
				int(c[1].X)+x, int(c[1].Y*25)+2+yVerse,
				"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
			)
		}
		canv.Gend()
		canv.Gend()

		y += 35

	}

	canv.Gend()
	return VerseInfo{
		MarginBottom: y,
	}
}
