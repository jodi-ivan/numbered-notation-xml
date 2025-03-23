package lyric

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type versePosition struct {
	Col      int
	Row      int
	RowWidth int
	Style    VerseRowStyle
}

// TODO: add the lyric notes like on kj-5, verse 2
func (li *lyricInteractor) RenderVerse(ctx context.Context, canv canvas.Canvas, y int, verses []repository.HymnVerse) VerseInfo {
	canv.Group("class='verses'", "style='font-family:Caladea'")

	allVerse := map[int][]string{}
	lineLength := float64(0)
	maxRightPost := float64(0)
	allCombine := map[int][][2]entity.Coordinate{}
	versePos := map[int]versePosition{}
	marginRight := 0
	multiColumn := false
	yPosRow := map[int]int{}
	maxY := float64(0)

	for _, verse := range verses {

		combine := [][2]entity.Coordinate{}
		whole := [][]LyricWordVerse{}

		err := json.Unmarshal([]byte(verse.Content.String), &whole)
		if err != nil {
			log.Println("[RenderVerse] failed to unmarshal, err ", err)
		}
		yPosRow[int(verse.Row.Int16)] = y + (25 * len(whole) * (int(verse.Row.Int16) - 1)) + ((int(verse.Row.Int16) - 1) * 35)

		style := VerseRowStyle(verse.StyleRow.Int32)

		if int(style) == 0 {
			style = VerseRowStyleSingleColumn
		}

		versePos[int(verse.VerseNum.Int32)] = versePosition{
			Col:      int(verse.Col.Int16),
			RowWidth: int(verse.StyleRow.Int32),
			Row:      int(verse.Row.Int16),
			Style:    style,
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
								x2 = x1 + li.CalculateLyricWidth(v.Text)
								break
							} else {
								x1 += li.CalculateLyricWidth(v.Text)

							}
						}

						startPosition := li.CalculateLyricWidth(lineText) + li.CalculateLyricWidth(wordPart)
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
			lineLength = math.Max(lineLength, li.CalculateLyricWidth(lineText))
			if verse.Col.Int16 == 1 {
				marginRight = int(lineLength)
			} else if verse.Col.Int16 == 2 {
				maxRightPost = math.Max(maxRightPost, lineLength)
				multiColumn = multiColumn || true
			}
			blob = append(blob, lineText)
		}
		allVerse[int(verse.VerseNum.Int32)] = blob
		allCombine[int(verse.VerseNum.Int32)] = combine
	}

	defaultX := int(math.Round((constant.LAYOUT_WIDTH / 2) - (lineLength / 2)))
	x := defaultX
	if multiColumn {
		// x = int(math.Round(constant.LAYOUT_WIDTH/2)) - ((int(lineLength) + int(maxRightPost)) / 2)
		x = constant.LAYOUT_INDENT_LENGTH * 2

	}
	totalVerse := len(allVerse)

	for i := 1; i < totalVerse+1; i++ {

		canv.Group("class='verse'")
		yVerse := y

		currentVerse := allVerse[i+1]

		// number verse
		margin := 0
		if versePos[i+1].Col == 2 && versePos[i+1].RowWidth == 6 {
			margin = constant.LAYOUT_WIDTH - (constant.LAYOUT_INDENT_LENGTH * 3) - int(maxRightPost)
			yVerse = yPosRow[versePos[i+1].Row]
			y = yVerse
			_ = marginRight
		}

		if versePos[i+1].Style == VerseRowStyleSingleColumn {
			x = defaultX + int(constant.LAYOUT_INDENT_LENGTH/2)
		}

		canv.Text(x-5-int(li.CalculateLyricWidth(fmt.Sprintf("%d. ", i+1)))+margin, y, fmt.Sprintf("%d. ", i+1))
		for _, liveVerse := range currentVerse {
			canv.Text(x+margin, y, liveVerse)
			y += 25

		}
		canv.Group()

		for _, c := range allCombine[i+1] {
			canv.Qbez(
				int(c[0].X)+x+margin, int(c[0].Y*25)+2+yVerse,
				int(c[0].X)+x+margin+(int(c[1].X-c[0].X)/2), yVerse+7+(int(c[0].Y)*25),
				int(c[1].X)+x+margin, int(c[1].Y*25)+2+yVerse,
				"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
			)
		}
		canv.Gend()
		canv.Gend()

		y += 35
		maxY = math.Max(maxY, float64(y))

	}

	canv.Gend()
	return VerseInfo{
		MarginBottom: int(maxY),
	}
}
