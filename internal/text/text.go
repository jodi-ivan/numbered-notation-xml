package text

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/lyric"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
	"github.com/jodi-ivan/numbered-notation-xml/utils/params"
)

type Text interface {
	NoteHasText(measureText []musicxml.MeasureText, t ...string) bool
	MeasureHasText(measure musicxml.Measure, t string) bool
	SetMeasureTextRenderer(ctx context.Context, noteRenderer *entity.NoteRenderer, note musicxml.Note, isLastNote bool) bool
	RenderMeasureText(ctx context.Context, y int, canv canvas.Canvas, notes []*entity.NoteRenderer, linestaff ...lines.LineStaff)
}

func NewText(l lyric.Lyric) Text {
	return &textInteractor{
		Lyric: l,
	}
}

type textInteractor struct {
	Lyric lyric.Lyric
}

func (ti *textInteractor) NoteHasText(measureText []musicxml.MeasureText, t ...string) bool {
	if len(t) == 0 {
		return len(measureText) > 0
	}

	matchCount := 0

	strorage := map[string]bool{}

	for _, nt := range measureText {
		strorage[nt.Text] = true
	}

	for _, it := range t {
		if strorage[it] {
			matchCount++
		}
	}

	return matchCount == len(t)
}

func (ti *textInteractor) MeasureHasText(measure musicxml.Measure, t string) bool {
	return measure.RightMeasureText != nil && measure.RightMeasureText.Text == t
}

func (ti *textInteractor) SetMeasureTextRenderer(ctx context.Context, noteRenderer *entity.NoteRenderer, note musicxml.Note, isLastNote bool) bool {
	affectMarginBottom := []string{DEFAULT_TEXT_REFREIN, DEFAULT_TEXT_FINE}
	prm, _ := params.GetParamFromContext(ctx)
	count := 0
	for _, mt := range note.MeasureText {
		if noteRenderer.MeasureText == nil {
			noteRenderer.MeasureText = []musicxml.MeasureText{}
		}

		if prm.Verse != 0 && strings.Contains(strings.ToLower(mt.Text), fmt.Sprintf("bait %d", prm.Verse)) {
			continue
		}
		alignment := musicxml.TextAlignmentLeft
		if isLastNote {
			alignment = musicxml.TextAlignmentRight
		}
		if slices.Contains(affectMarginBottom, mt.Text) {
			count++
		}

		noteRenderer.MeasureText = append(noteRenderer.MeasureText, musicxml.MeasureText{
			Text:      mt.Text,
			RelativeY: mt.RelativeY, TextAlignment: alignment,
		})
	}

	return count > 0

}

func (ti *textInteractor) RenderMeasureText(ctx context.Context, y int, canv canvas.Canvas, notes []*entity.NoteRenderer, linestaff ...lines.LineStaff) {

	for notePos, note := range notes {

		noteMarginTop := 0.0

		if len(note.MeasureText) > 0 {
			sort.Slice(note.MeasureText, func(i, j int) bool {
				return note.MeasureText[i].RelativeY < note.MeasureText[j].RelativeY
			})

			offset := 0
			if note.Fermata != nil {
				offset = FERMATA_OFFSET
			}

			if notes[0].Barline != nil && notes[0].Barline.Ending != nil {
				offset += REPEAT_LINE_OFFSET
			}

			if len(linestaff) > 0 {
				noteMarginTop = GetTextMarginBottom(linestaff[0], notes, notePos)
				offset += int(noteMarginTop)
			}

			for i, t := range note.MeasureText {

				style := []string{`font-style:italic`}
				if t.Text != DEFAULT_TEXT_REFREIN && t.Text != DEFAULT_TEXT_FINE {
					style = append(style, `font-size:65%`, `font-weight:bold`)
				}
				xPos := note.PositionX
				if t.TextAlignment == musicxml.TextAlignmentRight {
					textLength := ti.Lyric.CalculateLyricWidth(t.Text)
					xPos = constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH - int(textLength)
				}

				origPos := (len(note.MeasureText) - 1) * TEXT_BASELINE_DISTANCE
				yPos := (y - origPos) - offset - TEXT_TO_STAFF_DISTANCE - (i * -TEXT_BASELINE_DISTANCE)
				if t.RelativeY < 0 {
					yPos = y + (i * TEXT_BASELINE_DISTANCE) + (len(note.Lyric) * TEXT_TO_STAFF_DISTANCE) + 20
					style = []string{"font-size:60%", "font-family:'Figtree'", "font-weight:600"}
				}
				canv.Text(xPos, yPos, t.Text, fmt.Sprintf(`style="%s"`, strings.Join(style, ";")))
			}
		}
	}

}
