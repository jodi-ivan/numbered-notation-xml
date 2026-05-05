package barline

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderGregorian(canv canvas.Canvas, b *musicxml.Barline, isLastNote bool, lines [5]int, pos entity.Coordinate) {
	switch b.BarStyle {
	case musicxml.BarLineStyleRegular:
		xPos := pos.X
		if isLastNote {
			xPos += 4
		}
		canv.Line(int(xPos), lines[0], int(xPos), lines[4], "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9")

	case musicxml.BarLineStyleLightLight:
		canv.Line(int(pos.X)+1, lines[0], int(pos.X)+1, lines[4], "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9")
		canv.Line(int(pos.X)+4, lines[0], int(pos.X)+4, lines[4], "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9")

	case musicxml.BarLineStyleLightHeavy:
		canv.Line(int(pos.X)+1, lines[0], int(pos.X)+1, lines[4], "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9")
		canv.Line(int(pos.X)+6, lines[0]+2, int(pos.X)+6, lines[4]-2, "fill:none;stroke:#000000;stroke-linecap:square;stroke-width:4.6")

		if b.Repeat != nil && b.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward {
			canv.Text(int(pos.X)-6, lines[2]+5, ":", `style="font-family:Noto Music;font-size:0.6em"`)
		}

	case musicxml.BarLineStyleHeavyLight:
		canv.Line(int(pos.X)+2, lines[0]+2, int(pos.X)+2, lines[4]-2, "fill:none;stroke:#000000;stroke-linecap:square;stroke-width:4.6")
		canv.Line(int(pos.X)+7, lines[0], int(pos.X)+7, lines[4], "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.9")

		if b.Repeat != nil && b.Repeat.Direction == musicxml.BarLineRepeatDirectionForward {
			canv.Text(int(pos.X)+8, lines[2]+5, ":", `style="font-family:Noto Music;font-size:0.6em"`)
		}

	}
}

//
