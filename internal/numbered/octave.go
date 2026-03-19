package numbered

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func (ni *numberedInteractor) RenderOctave(ctx context.Context, canv canvas.Canvas, octave int, pos entity.Coordinate) {
	switch octave {
	case 1:
		canv.Circle(int(pos.X)+5, int(pos.Y)-15, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.75")
	case -1:
		canv.Circle(int(pos.X)+5, int(pos.Y)+5, 1, "fill:#000000;fill-opacity:1;stroke:#000000;stroke-width:0.75")
	}

}
