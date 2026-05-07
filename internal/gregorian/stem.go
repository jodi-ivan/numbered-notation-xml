package gregorian

import (
	"cmp"
	"math"
	"slices"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderStemUp(canv canvas.Canvas, lines [5]int, pos ...CoordinateWithNoteLength) (float64, float64) {
	start, end := []CoordinateWithNoteLength{}, []CoordinateWithNoteLength{}
	for _, v := range pos {
		x := float64(v.X) + 9
		y1, y2 := v.Y, v.Y-24
		start = append(start, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y1), NoteLength: v.NoteLength, Beam: v.Beam})
		end = append(end, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y2), NoteLength: v.NoteLength, Beam: v.Beam})
	}

	return renderStem(canv, lines, 1, start, end)
}

func RenderStemDown(canv canvas.Canvas, lines [5]int, pos ...CoordinateWithNoteLength) (float64, float64) {
	start, end := []CoordinateWithNoteLength{}, []CoordinateWithNoteLength{}
	for _, v := range pos {
		x := float64(v.X) + 0.5
		y1, y2 := (v.Y + 2), (v.Y + 28)
		start = append(start, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y1), NoteLength: v.NoteLength, Beam: v.Beam})
		end = append(end, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y2), NoteLength: v.NoteLength, Beam: v.Beam})
	}

	return renderStem(canv, lines, -1, start, end)

}

func renderStem(canv canvas.Canvas, lines [5]int, direction int, start, end []CoordinateWithNoteLength) (additional, clamp float64) {

	// end to end the grouping
	x1, y1 := end[0].X, end[0].Y
	x2, y2 := start[len(start)-1].X, end[len(end)-1].Y

	diffY := y2 - y1
	maxRise := STAFF_SPACE_WIDTH * 2.0 // Clamp the slope to 2 staff spaces

	if math.Abs(diffY) > maxRise {
		if diffY > 0 {
			y2 = y1 + maxRise
		} else {
			y2 = y1 - maxRise
		}

		clamp = diffY - maxRise
	}

	minDistance := STAFF_SPACE_WIDTH * 2.5

	for i := 0; i < len(start); i++ {
		// calculate the Y for the stem reach the beam
		x3 := start[i].X
		y3 := end[i].Y
		if len(start) > 1 {
			y3 = y1 + (x3-x1)*((y2-y1)/(x2-x1))
		}

		if math.Abs(y3-start[i].Y) < minDistance {
			additional = -1 * float64(direction) * (math.Abs(minDistance-math.Abs(y3-start[i].Y)) * (STAFF_SPACE_WIDTH / 2))
		}
	}

	// offset the position when beam line is horizontal line AND layering staff line
	intersect := slices.Index([]int{lines[0], lines[1], lines[2], lines[3], lines[4]}, int(y2))
	if y1 == y2 && intersect >= 0 {
		additional += (-1 * float64(direction)) + 3.5
	}

	canv.Group(`class="stem"`)
	for i := 0; i < len(start); i++ {
		x3 := start[i].X
		y3 := end[i].Y
		if len(start) > 1 {
			y3 = y1 + (x3-x1)*((y2-y1)/(x2-x1))
		}

		intersect := slices.Index([]int{lines[0], lines[1], lines[2], lines[3], lines[4]}, int(y3))
		if intersect >= 0 && (start[i].NoteLength == musicxml.NoteLengthQuarter || start[i].NoteLength == musicxml.NoteLengthHalf) {
			y3 += 0.5 * float64(direction)
		}
		canv.LineFloat64(start[i].X, start[i].Y, end[i].X, y3+additional, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
	}
	canv.Gend()

	return additional, clamp

}

func RenderGroupBeam(canv canvas.Canvas, groupBeam []CoordinateWithNoteLength, lines [5]int) {

	canv.Group(`class="beam-group"`)
	startPos, endPos := groupBeam[0], groupBeam[len(groupBeam)-1]

	farthestRank := slices.Clone(groupBeam)

	slices.SortFunc(farthestRank, func(a, b CoordinateWithNoteLength) int {
		return cmp.Compare(math.Abs(a.Y-float64(lines[2])), math.Abs(b.Y-float64(lines[2])))
	})
	farthest := farthestRank[len(farthestRank)-1]
	compared := cmp.Compare(farthest.Y, float64(lines[2]))

	y2Pos := endPos.Y
	stemOffset, clamp := renderMap[compared](canv, lines, groupBeam...)

	if y2Pos+(-1*float64(compared)*clamp) < startPos.Y && compared == 1 { // FIXME: y2 at min distance and stretch the y1 to satify slope
		/*
			### The Priority List (Top to Bottom)
			1.  **Min Distance (Critical):** No stem can ever be shorter than $\approx 2.5$ spaces. If it is, the whole beam must move.
			2.  **Max Slope (Limit):** The vertical difference between $y_1$ and $y_2$ should not exceed $2.0$ spaces.
			3.  **Standard Length (Goal):** Stems should ideally be $3.0$ or $3.5$ spaces, but this is the first thing you sacrifice to satisfy the two rules above.

		*/
		y2Pos += -1 * float64(compared) * clamp
	}

	if len(groupBeam) == 1 {
		offset := map[int]entity.Coordinate{
			-1: entity.NewCoordinate(0, +27),
			0:  entity.NewCoordinate(0, +27),
			1:  entity.NewCoordinate(9, -23),
		}

		canv.TextUnescaped(groupBeam[0].X+offset[compared].X,
			groupBeam[0].Y+offset[compared].Y,
			singleFlagHex[compared][groupBeam[0].NoteLength])

		canv.Gend()
		return
	}

	// BIG BEAM FLAG
	if compared <= 0 { // down.
		canv.LineFloat64(startPos.X+0.5, startPos.Y+stemOffset+27, endPos.X+0.5, y2Pos+stemOffset+27, `style="fill:none;stroke:#000000;stroke-linecap:butt;stroke-width:3"`)
	} else {
		canv.LineFloat64(startPos.X+9, startPos.Y+stemOffset-23, endPos.X+9, y2Pos+stemOffset-23, `style="fill:none;stroke:#000000;stroke-linecap:butt;stroke-width:3"`)
	}

	offsets := map[int][2]float64{
		-1: {0.5, 27},
		0:  {0.5, 27},
		1:  {9, -23},
	}

	direction := map[int]float64{
		-1: -5,
		0:  -5,
		1:  5,
	}

	total16 := 0
	pair := [][2]CoordinateWithNoteLength{{}}
	for _, v := range groupBeam {
		if len(v.Beam) > 1 {
			total16++
		}
		b := v.Beam[2]

		switch b.Type {
		case musicxml.NoteBeamTypeBegin:
			pair = append(pair, [2]CoordinateWithNoteLength{v})
			continue
		case musicxml.NoteBeamTypeEnd:
			currPair := pair[len(pair)-1]
			currPair[1] = v
			pair[len(pair)-1] = currPair
		}
	}

	if total16 > 0 {

		xOg1, yOg1 := startPos.X+offsets[compared][0], startPos.Y+stemOffset+offsets[compared][1]
		xOg2, yOg2 := endPos.X+offsets[compared][0], y2Pos+stemOffset+offsets[compared][1]

		for _, p := range pair {

			x1 := p[0].X + offsets[compared][0]
			x2 := p[1].X + offsets[compared][0]

			y1 := yOg1 + (x1-xOg1)*((yOg2-yOg1)/(xOg2-xOg1))
			y2 := yOg1 + (x2-xOg1)*((yOg2-yOg1)/(xOg2-xOg1))

			if p[0].IsEmpty() && !p[1].IsEmpty() {
				x3 := p[1].X - (constant.LOWERCASE_LENGTH / 2) + offsets[compared][0]
				y3 := y1 + (x3-x1)*((y2-y1)/(x2-x1))
				x1, y1 = x3, y3
			} else if !p[0].IsEmpty() && p[1].IsEmpty() {
				x3 := p[0].X + (constant.LOWERCASE_LENGTH / 2) + offsets[compared][0]
				y3 := y1 + (x3-x1)*((y2-y1)/(x2-x1))
				x2, y2 = x3, y3
			}

			canv.LineFloat64(x1, y1+direction[compared], x2, y2+direction[compared], `style="fill:none;stroke:#000000;stroke-linecap:butt;stroke-width:3"`)

		}

	}

	canv.Gend()

}
