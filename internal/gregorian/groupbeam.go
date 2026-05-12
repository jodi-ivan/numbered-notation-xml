package gregorian

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderGroupBeam(canv canvas.Canvas, lineStaff lines.LineStaff, groupBeam [][]CoordinateWithNoteLength, groupBeamSlurTies []SlurTieGroup) (map[int][]CoordinateWithNoteLength, VMargin) {
	directions := map[int][]CoordinateWithNoteLength{
		1: {}, -1: {},
	}

	margin := VMargin{
		Top:           entity.NewCoordinate(0, float64(lineStaff.GetTopLine())),
		Bottom:        entity.NewCoordinate(0, float64(lineStaff.GetBottomLine())),
		DefaultTop:    lineStaff.GetTopLine(),
		DefaultBottom: lineStaff.GetBottomLine(),
	}

	canv.Group(`class="beam-groups"`)
	for i, gr := range groupBeam {
		if len(gr) == 0 {
			continue
		}
		gMargin, direction := renderGroupBeam(canv, gr, lineStaff, groupBeamSlurTies)
		margin.Merge(gMargin)
		directions[direction] = append(directions[i], gr...)
	}
	canv.Gend()

	return directions, margin

}
