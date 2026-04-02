package verse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/footnote"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/svc/repository"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

var li = lyric.NewLyric()

func RenderVerse(ctx context.Context, canv canvas.Canvas, y int, verses map[int]repository.HymnVerse, verseFootnote map[int]map[int]repository.VerseFootNotes) VerseInfo {
	canv.Group("class='verses'", "style='font-family:Caladea'")

	allVerse := map[int][]string{}
	lineLength := float64(0)
	maxRightPost := float64(0)
	allCombine := map[int][][2]entity.Coordinate{}
	versePos := map[int]versePosition{}
	multiColumn := false
	yPosRow := map[int]int{}
	maxY := float64(0)

	for i := 2; i <= len(verses)+1; i++ {
		verse := verses[i]
		combine := [][2]entity.Coordinate{}
		whole := [][]LyricWordVerse{}

		err := json.Unmarshal([]byte(verse.Content.String), &whole)
		if err != nil {
			log.Println("[RenderVerse] failed to unmarshal, err ", err)
		}
		yPosRow[int(verse.Row.Int16)] = y + (25 * len(whole) * (int(verse.Row.Int16) - 1)) + ((int(verse.Row.Int16) - 1) * VERSE_SEPARATOR)

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
								if len(v.Text) > 0 {
									x2 -= (li.CalculateLyricWidth(string(v.Text[len(v.Text)-1])) / 2)
								}

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
				// marginRight = int(lineLength)
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

		canv.Group("class='verse'", fmt.Sprintf("number='%d'", i+1))
		yVerse := y

		currentVerse := allVerse[i+1]

		// number verse
		margin := 0
		if versePos[i+1].Col == 2 && versePos[i+1].RowWidth == 6 {
			margin = constant.LAYOUT_WIDTH - (constant.LAYOUT_INDENT_LENGTH * 3) - int(maxRightPost)
			yVerse = yPosRow[versePos[i+1].Row]
			y = yVerse
		}

		if versePos[i+1].Style == VerseRowStyleSingleColumn {
			x = defaultX + int(constant.LAYOUT_INDENT_LENGTH/2)
		}

		canv.Text(x-5-int(li.CalculateLyricWidth(fmt.Sprintf("%d. ", i+1)))+margin, y, fmt.Sprintf("%d. ", i+1))
		for line, liveVerse := range currentVerse {
			canv.Text(x+margin, y, liveVerse)

			if footnotes, hasFootnotes := verseFootnote[i+1]; hasFootnotes {
				currentLine, lineHasFootnotes := footnotes[line+1]

				verseStyle := footnote.VerseNoteStyle(currentLine.MarkerStyle.Int32)
				if lineHasFootnotes && verseStyle != footnote.VerseNoteStyleHeadless {
					xPos := x + margin + int(li.CalculateLyricWidth(liveVerse))
					styleFontSize := "font-family:'Figtree';font-weight:600;"
					switch verseStyle {
					case footnote.VerseNoteStyleAlignRight:
						styleFontSize = "font-family:'Figtree';font-size:60%;font-weight:600;"
						approxLineLength := constant.LAYOUT_WIDTH - (2 * defaultX)
						xPos = x + margin + approxLineLength
					case footnote.VerseNoteStyleHeadonly:
						styleFontSize = "font-family:'Caladea';font-size:90%;font-weight:600;"
						xPos -= int(li.CalculateLyricWidth(" ")) // lyric on db is just white spaces
					}
					canv.Group("class='footnotes'", fmt.Sprintf(`style="font-style:italic;%s"`, styleFontSize))
					canv.Text(xPos, y, currentLine.FootnoteMarker.String)
					canv.Gend()
				}
			}

			y += 25

		}

		if len(allCombine) > 0 {
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
		}
		canv.Gend()

		y += VERSE_SEPARATOR
		maxY = math.Max(maxY, float64(y))

	}

	canv.Gend()
	return VerseInfo{
		MarginBottom: int(maxY),
	}
}
