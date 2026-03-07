package lyric

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (li *lyricInteractor) RenderElision(ctx context.Context, canv canvas.Canvas, text []entity.Text, lyricPart int, pos entity.Coordinate) {
	offsetLyric := ""
	yPos := int(pos.Y) + 2
	for _, t := range text {

		if t.Underline == 1 {
			currTextLength := li.CalculateLyricWidth(t.Value)
			offset := li.CalculateLyricWidth(offsetLyric)
			canv.Qbez(
				int(pos.X+offset), yPos,
				int(pos.X+offset+(currTextLength/2)), yPos+6,
				int(pos.X+offset+currTextLength), yPos,
				"fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.1",
			)
		} else {
			offsetLyric += t.Value
		}
	}
}
