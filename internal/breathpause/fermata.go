package breathpause

import (
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderFermata(ctx context.Context, canv canvas.Canvas, fermata *musicxml.Femata, pos entity.Coordinate) {
	if fermata != nil {
		fermataUnicode := `&#x1D110;`

		fmt.Fprintf(
			canv.Writer(),
			`<text x="%.3f" y="%.3f" style="font-family:Noto Music;font-size:200%%"> %s </text>`,
			pos.X-5.5, pos.Y-5, fermataUnicode,
		)
	}
}
