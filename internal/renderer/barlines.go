package renderer

import (
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

// only support noto-music font
// FIXED: Print barline it as glyph
var unicode = map[musicxml.BarLineStyle]string{
	musicxml.BarLineStyleRegular:    `&#x01D100;`,
	musicxml.BarLineStyleLightHeavy: `&#x01D102;`,
	musicxml.BarLineStyleLightLight: `&#x01D101;`,
	musicxml.BarLineStyleHeavyHeavy: `&#x01D101;`,
	musicxml.BarLineStyleHeavyLight: `&#x01D103;`,
}

func RenderBarline(ctx context.Context, canv canvas.Canvas, barline musicxml.Barline, coordinate entity.Coordinate) {
	forward := ""
	backward := ""

	if barline.Repeat != nil {
		if barline.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward {
			backward = fmt.Sprintf(`<tspan x="%f" y="%f">:</tspan>`, coordinate.X-5, coordinate.Y-1)
		} else if barline.Repeat.Direction == musicxml.BarLineRepeatDirectionForward {
			//FIXED: adjust the size and position of forward barline
			forward = fmt.Sprintf(`<tspan x="%f" y="%f">:</tspan>`, coordinate.X+10, coordinate.Y-1)
		}
	}
	fmt.Fprintf(canv.Writer(), `<text x="%f" y="%f" style="font-family:Noto Music">
		%s
		<tspan x="%f" y="%f" font-size="130%%"> %s </tspan>
		%s
		</text>`,
		coordinate.X,
		coordinate.Y,
		backward,
		coordinate.X,
		coordinate.Y+3, unicode[barline.BarStyle], forward)
}
