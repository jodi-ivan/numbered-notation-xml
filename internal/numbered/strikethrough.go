package numbered

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (ni *numberedInteractor) RenderStrikethrough(ctx context.Context, canv canvas.Canvas, strikethrough bool, pos entity.Coordinate) {
	if strikethrough {
		canv.Line(int(pos.X)+10, int(pos.Y)-16, int(pos.X), int(pos.Y)+5, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.45")
	}
}
