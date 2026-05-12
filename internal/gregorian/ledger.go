package gregorian

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderLedgerLine(canv canvas.Canvas, pos entity.Coordinate, topLine, bottomLine int) {
	if pos.Y-float64(bottomLine) >= STAFF_SPACE_WIDTH {
		// ledger lines

		canv.Group(`class="ledger-lines"`)
		for ledgerPos := bottomLine + STAFF_SPACE_WIDTH; ledgerPos <= int(pos.Y); ledgerPos += STAFF_SPACE_WIDTH {
			x1 := pos.X - (constant.LOWERCASE_LENGTH / 2) + 3
			x2 := pos.X + 6 + (constant.LOWERCASE_LENGTH / 2)
			canv.Line(int(x1), ledgerPos, int(x2), ledgerPos, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.8")

		}
		canv.Gend()
	} else if float64(topLine)-pos.Y >= STAFF_SPACE_WIDTH {
		canv.Group(`class="ledger-lines"`)

		for ledgerPos := topLine - STAFF_SPACE_WIDTH; ledgerPos >= int(pos.Y); ledgerPos -= STAFF_SPACE_WIDTH {
			x1 := pos.X - (constant.LOWERCASE_LENGTH / 2) + 3
			x2 := pos.X + 6 + (constant.LOWERCASE_LENGTH / 2)
			canv.Line(int(x1), ledgerPos, int(x2), ledgerPos, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.8")

		}
		canv.Gend()

	}
}
