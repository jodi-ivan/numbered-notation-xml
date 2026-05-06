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
		start = append(start, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y1), NoteLength: v.NoteLength, Beam: v.Beam})
		end = append(end, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y2), NoteLength: v.NoteLength, Beam: v.Beam})
	}

	renderStem(canv, lines, 1, start, end)
}

func RenderStemDown(canv canvas.Canvas, lines [5]int, pos ...CoordinateWithNoteLength) {
	start, end := []CoordinateWithNoteLength{}, []CoordinateWithNoteLength{}
	for _, v := range pos {
		x := float64(v.X) + 0.5
		y1, y2 := (v.Y + 2), (v.Y + 28)
		start = append(start, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y1), NoteLength: v.NoteLength, Beam: v.Beam})
		end = append(end, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(x, y2), NoteLength: v.NoteLength, Beam: v.Beam})
	}

	renderStem(canv, lines, -1, start, end)

}

func renderStem(canv canvas.Canvas, lines [5]int, direction int, start, end []CoordinateWithNoteLength) {

	x1, y1 := end[0].X, end[0].Y
	x2, y2 := start[len(start)-1].X, end[len(end)-1].Y

	for i := 0; i < len(start); i++ {
		x3 := start[i].X
		y3 := end[i].Y
		if len(start) > 1 {
			y3 = y1 + (x3-x1)*((y2-y1)/(x2-x1))
		}

		intersect := slices.Index([]int{lines[0], lines[1], lines[2], lines[3], lines[4]}, int(y2))
		if intersect >= 0 && (start[i].NoteLength == musicxml.NoteLengthQuarter || start[i].NoteLength == musicxml.NoteLengthHalf) {
			y3 += 2.5 * float64(direction)
		}
		canv.LineFloat64(start[i].X, start[i].Y, end[i].X, y3, `style="fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1"`)
	}

}

func RenderGroupBeam(canv canvas.Canvas, groupBeam []CoordinateWithNoteLength, lines [5]int) {

	startPos, endPos := groupBeam[0], groupBeam[len(groupBeam)-1]

	farthestRank := slices.Clone(groupBeam)

	slices.SortFunc(farthestRank, func(a, b CoordinateWithNoteLength) int {
		return cmp.Compare(math.Abs(a.Y-float64(lines[2])), math.Abs(b.Y-float64(lines[2])))
	})
	farthest := farthestRank[len(farthestRank)-1]
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

	if endPos.NoteLength == musicxml.NoteLength16th && groupBeam[len(groupBeam)-2].NoteLength != endPos.NoteLength {
		// backward hook
		// TODO: forward hook? middle of group? why only handle the end of the group
		x1, y1 := startPos.X+offsets[compared][0], startPos.Y+offsets[compared][1]
		x2, y2 := endPos.X+offsets[compared][0], endPos.Y+offsets[compared][1]

		t := (1 - (0.7 * (1 / float64(len(groupBeam)))))

		mx := x1 + t*(x2-x1)
		my := y1 + t*(y2-y1)

		canv.LineFloat64(mx, my+direction[compared], x2, y2+direction[compared], `style="fill:none;stroke:#000000;stroke-linecap:butt;stroke-width:3"`)

	}

	total16 := 0
	pair := [][2]CoordinateWithNoteLength{}
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

	if total16%2 == 0 {

		xOg1, yOg1 := startPos.X+offsets[compared][0], startPos.Y+offsets[compared][1]
		xOg2, yOg2 := endPos.X+offsets[compared][0], endPos.Y+offsets[compared][1]

		for _, p := range pair {

			x1 := p[0].X + offsets[compared][0]
			x2 := p[1].X + offsets[compared][0]

			y1 := yOg1 + (x1-xOg1)*((yOg2-yOg1)/(xOg2-xOg1))
			y2 := yOg1 + (x2-xOg1)*((yOg2-yOg1)/(xOg2-xOg1))

			canv.LineFloat64(x1, y1+direction[compared], x2, y2+direction[compared], `style="fill:none;stroke:#000000;stroke-linecap:butt;stroke-width:3"`)

		}

	}

	// y3 = y1 + (x3-x1) * ((y2-y1)/(x2-x1))
}
