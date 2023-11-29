package renderer

import (
	"context"
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
