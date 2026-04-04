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

type Verse interface {
	RenderVerse(ctx context.Context, canv canvas.Canvas, y int, verses map[int]repository.HymnVerse, verseFootnote map[int]map[int]repository.VerseFootNotes) VerseInfo
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

func (v *verseInteractor) elisionPosition(p LyricPartVerse, y int, lineText, wordPart string) [2]entity.Coordinate {
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

	startPosition := v.Lyric.CalculateLyricWidth(lineText) + v.Lyric.CalculateLyricWidth(wordPart)
	return [2]entity.Coordinate{
		entity.NewCoordinate(startPosition+x1, float64(y)),
		entity.NewCoordinate(startPosition+x2, float64(y)),
	}
}

func (v *verseInteractor) parse(y int, verses map[int]repository.HymnVerse) ParsedVerseWithInfo {

	result := ParsedVerseWithInfo{
		Verses:        map[int]ParsedVerse{},
		IsMultiColumn: false,
		RowPositionY:  map[int]int{},
	}

	for i := 2; i <= len(verses)+1; i++ {
		verse := verses[i]
		whole := [][]LyricWordVerse{}

		err := json.Unmarshal([]byte(verse.Content.String), &whole)
		if err != nil {
			log.Println("[RenderVerse] failed to unmarshal, err ", err)
			return result
		}

		if _, ok := result.RowPositionY[int(verse.Row.Int16)]; !ok {
			result.RowPositionY[int(verse.Row.Int16)] = y + (LINE_DISTANCE * len(whole) * (int(verse.Row.Int16) - 1)) + ((int(verse.Row.Int16) - 1) * VERSE_SEPARATOR)
		}

		style := VerseRowStyle(verse.StyleRow.Int32)

		if int(style) == 0 {
			style = VerseRowStyleSingleColumn
		}
		parsedVerse := ParsedVerse{
			ElisionMarks: [][2]entity.Coordinate{},
			Position: versePosition{
				Col:      int(verse.Col.Int16),
				RowWidth: int(verse.StyleRow.Int32),
				Row:      int(verse.Row.Int16),
				Style:    style,
			},
		}

		for iLine, line := range whole {
			lineText := ""
			for _, word := range line {
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
			result.MaxLineWidth = math.Max(result.MaxLineWidth, v.Lyric.CalculateLyricWidth(lineText))
			if verse.Col.Int16 == 2 {
				result.MaxRightPos = math.Max(result.MaxRightPos, result.MaxLineWidth)
				result.IsMultiColumn = result.IsMultiColumn || true
			}
			parsedVerse.Verse = append(parsedVerse.Verse, lineText)
		}
		result.Verses[int(verse.VerseNum.Int32)] = parsedVerse
	}

	return result
}

func (v *verseInteractor) RenderVerse(ctx context.Context, canv canvas.Canvas, y int, verses map[int]repository.HymnVerse, verseFootnote map[int]map[int]repository.VerseFootNotes) VerseInfo {
	canv.Group("class='verses'", "style='font-family:Caladea'")

	parsedVerse := v.parse(y, verses)

	defaultX := int(math.Round((constant.LAYOUT_WIDTH / 2) - (parsedVerse.MaxLineWidth / 2)))
	x := defaultX
	if parsedVerse.IsMultiColumn {
		x = constant.LAYOUT_INDENT_LENGTH * 2
	}
	totalVerse := len(parsedVerse.Verses)

	maxY := float64(0)
	for i := 1; i < totalVerse+1; i++ {

		canv.Group("class='verse'", fmt.Sprintf("number='%d'", i+1))
		yVerse := y

		currentVerse := parsedVerse.Verses[i+1]

		// number verse
		margin := 0
		if currentVerse.Position.Col == 2 && VerseRowStyle(currentVerse.Position.RowWidth) == VerseRowStyleDualColumn {
			margin = constant.LAYOUT_WIDTH - (constant.LAYOUT_INDENT_LENGTH * 3) - int(parsedVerse.MaxRightPos)
			yVerse = parsedVerse.RowPositionY[currentVerse.Position.Row]
			y = yVerse
		}

		if currentVerse.Position.Style == VerseRowStyleSingleColumn {
			x = defaultX + int(constant.LAYOUT_INDENT_LENGTH/2)
		}

		prefixNum := fmt.Sprintf("%d. ", i+1)
		canv.Text(x-5-int(v.Lyric.CalculateLyricWidth(prefixNum))+margin, y, prefixNum)
		for line, liveVerse := range currentVerse.Verse {
			canv.Text(x+margin, y, liveVerse)
			cursor := footnote.VerseLineCursor{
				VerseNo:    i + 1,
				LinePos:    line + 1,
				Leftmargin: margin,
				LineText:   liveVerse,
			}
			v.Footnote.AssignFootnotesMarker(canv, entity.NewCoordinate(float64(x), float64(y)), defaultX, cursor, verseFootnote)
			y += LINE_DISTANCE
		}

		if len(currentVerse.ElisionMarks) > 0 {
			canv.Group()
			for _, c := range currentVerse.ElisionMarks {
				x0 := int(c[0].X) + x + margin
				x1 := int(c[1].X) + x + margin
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
