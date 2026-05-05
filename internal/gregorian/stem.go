package gregorian

import (
	"cmp"
	"math"
	"slices"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderStemUp(canv canvas.Canvas, lines [5]int, pos ...CoordinateWithNoteLength) {
	start, end := []CoordinateWithNoteLength{}, []CoordinateWithNoteLength{}
	for _, v := range pos {
		x := float64(v.X) + 9
		y1, y2 := v.Y, v.Y-24
		start = append(start, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y1), NoteLength: v.NoteLength})
		end = append(end, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y2), NoteLength: v.NoteLength})
	}

	renderStem(canv, lines, 1, start, end)
}

func RenderStemDown(canv canvas.Canvas, lines [5]int, pos ...CoordinateWithNoteLength) {
	start, end := []CoordinateWithNoteLength{}, []CoordinateWithNoteLength{}
	for _, v := range pos {
		x := float64(v.X) + 0.5
		y1, y2 := (v.Y + 2), (v.Y + 28)
		start = append(start, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y1), NoteLength: v.NoteLength})
		end = append(end, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y2), NoteLength: v.NoteLength})
	}

	renderStem(canv, lines, -1, start, end)

}

func renderStem(canv canvas.Canvas, lines [5]int, direction int, start, end []CoordinateWithNoteLength) {
	for i := 0; i < len(start); i++ {
		y2 := end[i].Y
		intersect := slices.Index([]int{lines[0], lines[1], lines[2], lines[3], lines[4]}, int(y2))
		if intersect >= 0 && (start[i].NoteLength == musicxml.NoteLengthQuarter || start[i].NoteLength == musicxml.NoteLengthHalf) {
			y2 += 2.5 * float64(direction)
		}
		canv.LineFloat64(start[i].X, start[i].Y, end[i].X, y2, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
	}

}

func RenderGroupBeam(canv canvas.Canvas, groupBeam []CoordinateWithNoteLength, lines [5]int) {

	startPos, endPos := groupBeam[0], groupBeam[len(groupBeam)-1]
	slices.SortFunc(groupBeam, func(a, b CoordinateWithNoteLength) int {
		return cmp.Compare(math.Abs(a.Y-float64(lines[2])), math.Abs(b.Y-float64(lines[2])))
	})
	farthest := groupBeam[len(groupBeam)-1]
	compared := cmp.Compare(farthest.Y, float64(lines[2]))
	renderMap[compared](canv, lines, groupBeam...)

	if len(groupBeam) == 1 {

		offset := map[int]entity.Coordinate{
			-1: entity.NewCoordinate(0, +27),
			0:  entity.NewCoordinate(0, +27),
			1:  entity.NewCoordinate(9, -23),
		}

		canv.TextUnescaped(groupBeam[0].X+offset[compared].X,
			groupBeam[0].Y+offset[compared].Y,
			singleFlagHex[compared][groupBeam[0].NoteLength])
		return
	}

	// BIG BEAM FLAG
	if compared <= 0 { // down.
		canv.LineFloat64(startPos.X+0.5, startPos.Y+27, endPos.X+0.5, endPos.Y+27, `style="fill:none;stroke:#000000;stroke-linecap:butt;stroke-width:3"`)
	} else {
		canv.LineFloat64(startPos.X+9, startPos.Y-23, endPos.X+9, endPos.Y-23, `style="fill:none;stroke:#000000;stroke-linecap:butt;stroke-width:3"`)

	}

}
