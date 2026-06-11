package lyric

import (
	"context"
	"strings"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (li *lyricInteractor) RenderElision(ctx context.Context, canv canvas.Canvas, text []entity.Text, lyricPart int, pos entity.Coordinate, style ...string) {
	offsetLyric := ""
	yPos := int(pos.Y) + 2
	for _, t := range text {

		if t.Underline == 1 {
			styleStr := "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1;"
			if len(style) > 0 {
				styleStr += strings.Join(style, ";")
			}
			currTextLength := li.CalculateLyricWidth(t.Value)
			offset := li.CalculateLyricWidth(offsetLyric)
			canv.Qbez(
				int(pos.X+offset), yPos,
				int(pos.X+offset+(currTextLength/2)), yPos+6,
				int(pos.X+offset+currTextLength), yPos,
				styleStr,
			)
		} else {
			offsetLyric += t.Value
		}
	}
}
