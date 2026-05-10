package gregorian

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"sort"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func GetYPosGroup(group []CoordinateWithNoteLength, lines [5]int) int {

	slices.SortFunc(group, func(a, b CoordinateWithNoteLength) int {
		return cmp.Compare(math.Abs(a.Y-float64(lines[2])), math.Abs(b.Y-float64(lines[2])))
	})
	farthest := group[len(group)-1]
	compared := cmp.Compare(farthest.Y, float64(lines[2]))
	if compared == 0 {
		compared = -1
	}

	return compared
}

func RenderSlurTies(canv canvas.Canvas, lineStaff LineStaff, groupBeam [][]CoordinateWithNoteLength, slurties []SlurTieGroup) int {
	canv.Group(`class="slurties"`)

	lines := lineStaff.GetLines()
	marginButtom := lines[4]

	tiesEndOffset := map[int][2]entity.Coordinate{
		-1: {
			entity.NewCoordinate(6, -6),
			entity.NewCoordinate(6, -6),
		},
		1: {
			entity.NewCoordinate(6, 6),
			entity.NewCoordinate(3, 6),
		},
	}

	tiesNotes := map[string]bool{}

	sort.Slice(slurties, func(i, j int) bool {
		return slurties[i].Ties != nil && slurties[j].Ties == nil
	})

	minMax := map[int]func(float64, float64) float64{
		-1: math.Min,
		1:  math.Max,
	}

	for _, st := range slurties {

		group := []CoordinateWithNoteLength{}

		for _, gs := range groupBeam {
			for _, g := range gs {
				if g.NoteID == st.NoteMember[0] {
					group = gs
					break
				}
			}
		}

		direction := st.AccumulativeDirection
		if direction < 0 {
			direction = -1
		} else {
			direction = 1
		}

		if len(group) > 0 && group[0].NoteID != st.NoteMember[0] {
			// you part of member, but not the root.
			// we just follow the group rules instead of consensus
			direction = GetYPosGroup(group, lines)
		}
		isTies := false
		var sluLineType musicxml.NoteSlurLineType

		if st.Ties != nil {
			tiesNotes[st.NoteMember[0]] = true
			isTies = true
			sluLineType = st.Ties.LineType
		} else if st.Slur != nil {
			sluLineType = st.Slur.LineType
		}
		// 3 cases.
		start := st.Start
		end := st.End

		placeholderY := float64(lines[2])
		if start.Y != 0 {
			placeholderY = start.Y
		} else if end.Y != 0 {
			placeholderY = end.Y
		}

		// no start
		if start.X == 0 || start.Y == 0 {
			start = entity.NewCoordinate(float64(lineStaff.GetLeftIndent())-(PADDING_WIDTH*2), placeholderY)
		}

		// no end
		if end.X == 0 || end.Y == 0 {
			end = entity.NewCoordinate(float64(lineStaff.GetMarginRight()), placeholderY)
		}

		getYVal := map[int]float64{
			-1: st.MinY,
			1:  st.MaxY,
		}
		// direction := st.AccumulativeDirection
		// if direction < 0 {
		// 	direction = -1
		// } else {
		// 	direction = 1
		// }
		pull := entity.NewCoordinate(
			start.X+((end.X-start.X)/2),
			getYVal[direction])

		if pull.Y <= 0 {
			pull.Y = minMax[direction](start.Y, end.Y) + (float64(direction))
		}

		pullY := 0.0

		block := ((end.X - start.X) / constant.UPPERCASE_LENGTH) // * 2
		if block < 2 {
			pullY += 15 * float64(direction)
		} else if block < 5 {
			pullY += 20 * float64(direction)
		} else {
			pullY += 25 * float64(direction) //long distance ties, need more height
		}

		startYOffset := 0
		if !isTies && tiesNotes[st.NoteMember[0]] {
			pullY += 10 * float64(direction)
			startYOffset = direction * 3
		}

		lineType := "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.5"
		if sluLineType == musicxml.NoteSlurLineTypeDashed {
			// Calculate approximate curve length
			curveLen := utils.QuadBezierLength(
				start,
				pull,
				end,
				30)

			dash := 3.5
			patternCount := math.Floor(curveLen / (dash * 2))
			gap := (curveLen / patternCount) - dash

			lineType += fmt.Sprintf(";stroke-dasharray:%.1f %.1f;", dash, gap)
			lineType += "stroke-dashoffset:" + fmt.Sprintf("%f", dash/2) + ";"

			canv.Qbez(
				int(math.Round(start.X))+int(tiesEndOffset[direction][0].X),
				int(math.Round(start.Y))+int(tiesEndOffset[direction][0].Y)+startYOffset,
				int(math.Round(pull.X))+int(tiesEndOffset[direction][0].X)/2,
				int(math.Ceil(pull.Y+pullY)),
				int(math.Round(end.X))+int(tiesEndOffset[direction][1].X),
				int(math.Round(end.Y))+int(tiesEndOffset[direction][1].Y),
				lineType,
				fmt.Sprintf(`direction="%d"`, direction),
			)
		} else {
			w := canv.Writer()
			fmt.Fprintf(w, `<path d="%s" style="stroke-width: 0.5;stroke: #000000;" direction="%d"></path>`,
				GenerateTaperedSlur(math.Round(start.X)+tiesEndOffset[direction][0].X,
					math.Round(start.Y)+tiesEndOffset[direction][0].Y+float64(startYOffset),
					math.Round(end.X)+tiesEndOffset[direction][1].X,
					math.Round(end.Y)+tiesEndOffset[direction][1].Y,
					math.Round(pull.X)+3+(tiesEndOffset[direction][0].X/2),
					math.Ceil(pull.Y+pullY),
				), direction)
		}

		if marginButtom < int(math.Ceil(pull.Y+pullY)) {
			marginButtom = int(math.Ceil(pull.Y + pullY))
		}

		canv.Circle(
			int(math.Round(pull.X))+3+int(tiesEndOffset[direction][0].X)/2,
			int(math.Ceil(pull.Y+pullY)), 2, "fill:none;stroke:#FF0000;stroke-linecap:round;stroke-width:0.4")
		// canv.Text(
		// 	(int(math.Round(pull.X))+int(tiesEndOffset[direction][0].X)/2)+3,
		// 	int(math.Ceil(pull.Y+pullY)),
		// 	fmt.Sprintf("Block: %.2f", block), "font-size:0.5em;font-family:Caladea")
	}
	canv.Gend()
	return marginButtom

}

func GenerateTaperedSlur(x1, y1, x2, y2, cx, cy float64) string {
	thickness := 4.0 // The "expensive" look thickness

	// Top control point (pulled higher)
	cyTop := cy - (thickness / 2)
	// Bottom control point (pulled lower)
	cyBottom := cy + (thickness / 2)

	return fmt.Sprintf(
		"M %f,%f Q %f,%f %f,%f Q %f,%f %f,%f Z",
		x1, y1, // Start point
		cx, cyTop, // Top "pull"
		x2, y2, // End point
		cx, cyBottom, // Bottom "pull" back
		x1, y1, // Back to start
	)
}
