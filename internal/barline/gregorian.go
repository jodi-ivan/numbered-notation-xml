package barline

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderGregorian(canv canvas.Canvas, b *musicxml.Barline, isLastNote bool, staffLine lines.LineStaff, pos entity.Coordinate) {
	canv.Group(`class="barline"`)
	defer func() {
		canv.Gend()
	}()
	topLine := staffLine.GetTopLine()
	bottomLine := staffLine.GetBottomLine()
	middleLine := staffLine.GetMiddleLine()

	switch b.BarStyle {
	case musicxml.BarLineStyleRegular:
		xPos := pos.X
		if isLastNote {
			xPos += 4
		}
		canv.Line(int(xPos), topLine, int(xPos), bottomLine, lineStyleRegular)

	case musicxml.BarLineStyleLightLight:
		leftBarlineOffset := 1
		rightBarlineOffset := 4
		canv.Line(
			int(pos.X)+leftBarlineOffset, topLine,
			int(pos.X)+leftBarlineOffset, bottomLine, lineStyleRegular)

		canv.Line(
			int(pos.X)+rightBarlineOffset, topLine,
			int(pos.X)+rightBarlineOffset, bottomLine, lineStyleRegular)

	case musicxml.BarLineStyleLightHeavy:
		canv.Line(
			int(pos.X)+1, topLine,
			int(pos.X)+1, bottomLine, lineStyleRegular)

		canv.Line(
			int(pos.X)+6, topLine+2, // 2pts  from round type offset
			int(pos.X)+6, bottomLine-2, lineStyleHeavy)

		if b.Repeat != nil && b.Repeat.Direction == musicxml.BarLineRepeatDirectionBackward {
			canv.Text(
				int(pos.X)-6, middleLine+5,
				":", `style="font-family:Noto Music;font-size:19.2px"`)
		}

	case musicxml.BarLineStyleHeavyLight:
		canv.Line(
			int(pos.X)+2, topLine+2,
			int(pos.X)+2, bottomLine-2, lineStyleHeavy)

		canv.Line(
			int(pos.X)+7, topLine,
			int(pos.X)+7, bottomLine, lineStyleRegular)

		if b.Repeat != nil && b.Repeat.Direction == musicxml.BarLineRepeatDirectionForward {
			canv.Text(
				int(pos.X)+8, middleLine+5,
				":", `style="font-family:Noto Music;font-size:19.2px"`)
		}

	}
}
