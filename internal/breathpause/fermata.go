package breathpause

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderFermata(ctx context.Context, canv canvas.Canvas, fermata *musicxml.Femata, pos entity.Coordinate) {
	if fermata != nil {
		fermataUnicode := `&#x1D110;`

		canv.TextUnescaped(
			pos.X-4, pos.Y-7,
			fermataUnicode,
			`style="font-family:Noto Music;font-size:150%"`,
		)
	}
}
