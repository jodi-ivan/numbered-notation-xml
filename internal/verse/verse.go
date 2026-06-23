package verse

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/footnote"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

type Verse interface {
	RenderVerse(ctx context.Context, canv canvas.Canvas, y int, metadata *entity.HymnMetaData) VerseInfo
}

type verseInteractor struct {
	Footnote footnote.Footnote
	Lyric    lyric.Lyric
}

func New(f footnote.Footnote, l lyric.Lyric) Verse {
	return &verseInteractor{
		Footnote: f,
		Lyric:    l,
	}
}

func (v *verseInteractor) elisionPosition(p entity.LyricPartVerse, y int, lineBeforeWord, syllableBeforeElision string) [2]entity.Coordinate {
	x1 := float64(0)
	x2 := float64(0)
	for _, b := range p.Breakdown {
		if b.Underline {
			x2 = x1 + v.Lyric.CalculateLyricWidth(b.Text)
			if len(b.Text) > 0 {
				x2 -= (v.Lyric.CalculateLyricWidth(string(b.Text[len(b.Text)-1])) / 2)
			}

			break
		} else {
			x1 += v.Lyric.CalculateLyricWidth(b.Text)

		}
	}

	startPosition := v.Lyric.CalculateLyricWidth(lineBeforeWord) + v.Lyric.CalculateLyricWidth(syllableBeforeElision)
	return [2]entity.Coordinate{
		entity.NewCoordinate(startPosition+x1, float64(y)),
		entity.NewCoordinate(startPosition+x2, float64(y)),
	}
}

func getPos(idx, style, totalVerse int) (row, col, newStyle int) {
	if (idx == 0 && totalVerse == 1) || style == 0 || style == 12 {
		return idx + 1, 1, 12
	}

	half := totalVerse / 2
	col = 1
	if (idx+1)%2 == 1 && idx+1 == totalVerse {
		col = 1
		row = half + 1 // Placed at the bottom across both columns
		style = 12
	} else {
		// 2. Vertical Column Logic
		if idx < half {
			// First Column (Vertical flow)
			col = 1
			row = idx + 1
		} else {
			// Second Column (Starts at Verse 4 if count is 5)
			col = 2
			row = (idx - half) + 1
		}
	}

	return row, col, style
}

func (v *verseInteractor) parse(ctx context.Context, y int, metadata *entity.HymnMetaData) ParsedVerseWithInfo {

	prm, _ := params.GetParamFromContext(ctx)

	result := ParsedVerseWithInfo{
		Verses:        map[int]ParsedVerse{},
		IsMultiColumn: false,
		RowPositionY:  map[int]int{},
	}

	if prm.Verse > 1 {
		delete(metadata.Verse, prm.Verse)
		delete(metadata.ParsedVerse, prm.Verse)

		if !prm.SingleVerseMode {
			delete(metadata.Verse, prm.Verse-1)
			delete(metadata.ParsedVerse, prm.Verse-1)
		}
	}

	if len(metadata.Verse) == 0 {
		return result
	}

	versesNo := utils.GetMapSortedKeys(metadata.Verse)

	for idx, i := range versesNo {

		verse := metadata.Verse[i]
		whole := metadata.ParsedVerse[i]

		row := int(verse.Row.Int16)
		col := int(verse.Col.Int16)
		style := int(verse.StyleRow.Int32)

		if prm.Verse > 1 {
			// get the 1st style instead, since there is cases when the verse is last verse
			// style can be modified.
			style = int(metadata.Verse[versesNo[0]].StyleRow.Int32)
			if style == 0 {
				style = int(VerseRowStyleSingleColumn)
			}
			row, col, style = getPos(idx, style, len(versesNo))
		}

		parsedVerse := ParsedVerse{
			ElisionMarks: [][2]entity.Coordinate{},
			Position: versePosition{
				Col: col, Row: row,
				Style: VerseRowStyle(style),
			},
		}

		totalLine := 0

		for iLine, line := range whole {
			lineText := ""
			for _, word := range line {
				if word.ScoreOnly {
					continue
				}
				wordPart := ""
				for _, p := range word.Breakdown {
					if p.Combine {
						elisionPos := v.elisionPosition(p, iLine, lineText, wordPart)
						parsedVerse.ElisionMarks = append(parsedVerse.ElisionMarks, elisionPos)
					}
					wordPart = wordPart + p.Text
				}
				lineText = lineText + " " + word.Word
			}
			if lineText == "" {
				continue
			}
			totalLine++
			result.MaxLineWidth = math.Max(result.MaxLineWidth, v.Lyric.CalculateLyricWidth(lineText))
			if col == 2 {
				result.MaxRightPos = math.Max(result.MaxRightPos, result.MaxLineWidth)
				result.IsMultiColumn = result.IsMultiColumn || true
			}
			parsedVerse.Verse = append(parsedVerse.Verse, lineText)
		}

		if _, ok := result.RowPositionY[row]; !ok {
			result.RowPositionY[row] = y + (LINE_DISTANCE * totalLine * (row - 1)) + ((row - 1) * VERSE_SEPARATOR)
		}
		result.Verses[i] = parsedVerse
	}

	return result
}

func (v *verseInteractor) RenderVerse(ctx context.Context, canv canvas.Canvas, y int, metadata *entity.HymnMetaData) VerseInfo {
	canv.Group("class='verses'", "style='font-family:Caladea;font-size:16px'")

	prm, _ := params.GetParamFromContext(ctx)

	parsedVerse := v.parse(ctx, y, metadata)

	defaultX := int(math.Round((constant.LAYOUT_WIDTH / 2) - (parsedVerse.MaxLineWidth / 2)))
	x := defaultX
	if parsedVerse.IsMultiColumn {
		x = constant.LAYOUT_INDENT_LENGTH * 2
	}
	totalVerse := len(parsedVerse.Verses)
	if prm.Verse != 0 {
		totalVerse -= 2
	}

	offset := 0.0

	maxY := float64(0)
	versesNo := utils.GetMapSortedKeys(parsedVerse.Verses)
	for _, i := range versesNo {

		canv.Group("class='verse'", fmt.Sprintf("number='%d'", i))
		yVerse := y

		currentVerse := parsedVerse.Verses[i]

		row := currentVerse.Position.Row
		col := currentVerse.Position.Col
		style := currentVerse.Position.Style

		// number verse
		margin := 0
		if parsedVerse.IsMultiColumn {
			if col == 2 {
				margin = constant.LAYOUT_WIDTH - (constant.LAYOUT_INDENT_LENGTH * 3.5) - int(parsedVerse.MaxRightPos)
				yVerse = parsedVerse.RowPositionY[row]
				y = yVerse
			} else if style == VerseRowStyleSingleColumn {
				margin = -1 * (constant.LAYOUT_INDENT_LENGTH / 4)
			}
		}

		if parsedVerse.IsMultiColumn && totalVerse > 3 && totalVerse%2 == 1 { // clamp the gap --> col 1 increase margin, col 2 decrease margin
			if parsedVerse.IsMultiColumn && constant.LAYOUT_WIDTH > parsedVerse.MaxLineWidth*4 {
				offset = parsedVerse.MaxLineWidth / 2
			}

			if col == 1 {
				offset = math.Abs(offset)
			} else {
				offset = math.Abs(offset) * -1
			}

			if style == VerseRowStyleSingleColumn {
				offset = (constant.LAYOUT_INDENT_LENGTH / 4)
			}
		}

		if style == VerseRowStyleSingleColumn {
			x = defaultX + int(constant.LAYOUT_INDENT_LENGTH/2)
		}
		xPos := x + margin + int(offset)

		prefixNum := fmt.Sprintf("%d. ", i)
		canv.Text(xPos-5-int(v.Lyric.CalculateLyricWidth(prefixNum)), y, prefixNum)
		for line, liveVerse := range currentVerse.Verse {
			if strings.HasPrefix(liveVerse, "    ") {
				liveVerse = strings.ReplaceAll(liveVerse, "    ", strings.Repeat("&#160;", 4))
			}
			canv.TextUnescaped(float64(xPos), float64(y), liveVerse)
			cursor := footnote.VerseLineCursor{
				VerseNo:    i,
				LinePos:    line + 1,
				Leftmargin: margin + int(offset),
				LineText:   liveVerse,
			}
			v.Footnote.AssignFootnotesMarker(canv, entity.NewCoordinate(float64(x), float64(y)), defaultX, cursor, metadata.VerseFootNotes)
			y += LINE_DISTANCE
		}

		if len(currentVerse.ElisionMarks) > 0 {
			canv.Group()
			for _, c := range currentVerse.ElisionMarks {
				x0 := int(c[0].X) + xPos
				x1 := int(c[1].X) + xPos
				xMid := x0 + (x1-x0)/2

				y0 := int(c[0].Y*LINE_DISTANCE) + ELISION_Y_OFFSET + yVerse
				y1 := int(c[1].Y*LINE_DISTANCE) + ELISION_Y_OFFSET + yVerse
				yCtrl := yVerse + ELISION_PULL_Y_OFFSET + int(c[0].Y)*LINE_DISTANCE

				canv.Qbez(x0, y0,
					xMid, yCtrl,
					x1, y1,
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
