package rhythm

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"sort"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/internal/utils"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func GetGroupSlueTies(notes []*entity.NoteRenderer, staffLine lines.LineStaff) []SlurTieGroup {
	var groups []SlurTieGroup
	var tieTracking *SlurTieGroup
	slurTracking := map[int]SlurTieGroup{}

	for _, note := range notes {
		yPos, direction := noteYPosDirection(note, staffLine)

		tieTracking = processTieNote(note, yPos, direction, tieTracking, &groups)

		if len(note.Slur) > 0 {
			for sid, slur := range note.Slur {
				processSlurNote(note, sid, slur, yPos, direction, staffLine, slurTracking, &groups)
			}
			continue
		}

		// middle note: accumulate into all active trackers
		accumulateMiddle(note.UUID, yPos, direction, tieTracking, slurTracking)
	}

	if tieTracking != nil {
		groups = append(groups, *tieTracking)
	}
	for _, v := range slurTracking {
		groups = append(groups, v)
	}
	return groups
}

func accumulateMiddle(uuid string, yPos float64, direction int, tie *SlurTieGroup, slurs map[int]SlurTieGroup) {
	if tie != nil && yPos != 0 {
		tie.MaxY = math.Max(tie.MaxY, yPos)
		tie.MinY = math.Min(tie.MinY, yPos)
		tie.AccumulativeDirection += direction
		tie.NoteMember = append(tie.NoteMember, uuid)
	}

	for i, v := range slurs {
		v.AccumulativeDirection += direction
		v.NoteMember = append(v.NoteMember, uuid)
		if yPos != 0 {
			v.MaxY = math.Max(v.MaxY, yPos)
			v.MinY = math.Min(v.MinY, yPos)
		}
		slurs[i] = v
	}
}

func newSlurTieGroup(yPos float64, staffLine lines.LineStaff) SlurTieGroup {
	if yPos == 0 {
		return SlurTieGroup{
			MaxY: float64(staffLine.GetTopLine()),
			MinY: float64(staffLine.GetBottomLine()),
		}
	}
	return SlurTieGroup{MaxY: yPos, MinY: yPos}
}

func noteYPosDirection(note *entity.NoteRenderer, staffLine lines.LineStaff) (float64, int) {
	if note.AbsoluteNote == "" || note.AbsoluteOctave <= 0 {
		return 0, 0
	}
	r := rune(note.AbsoluteNote[0])
	return staffLine.GetYPos(r, note.AbsoluteOctave), staffLine.GetStemDirection(r, note.AbsoluteOctave)
}

func processTieNote(
	note *entity.NoteRenderer,
	yPos float64,
	direction int,
	tracking *SlurTieGroup,
	groups *[]SlurTieGroup,
) *SlurTieGroup {
	if note.Tie == nil {
		return tracking
	}

	if tracking == nil && note.Tie.Type == musicxml.NoteSlurTypeStart && !note.Tie.NumberedOnly {
		return &SlurTieGroup{
			MaxY: yPos, MinY: yPos,

			Ties:  note.Tie,
			Start: entity.NewCoordinate(float64(note.PositionX), yPos),

			NoteMember:            []string{note.UUID},
			AccumulativeDirection: direction,
		}
	}

	if tracking != nil && note.Tie.Type == musicxml.NoteSlurTypeStop {
		tracking.AccumulativeDirection += direction
		tracking.End = entity.NewCoordinate(float64(note.PositionX), yPos)
		tracking.NoteMember = append(tracking.NoteMember, note.UUID)
		*groups = append(*groups, *tracking)
		return nil
	}

	return tracking
}
func processSlurNote(
	note *entity.NoteRenderer,
	sid int, slur entity.Slur,
	yPos float64, direction int,
	staffLine lines.LineStaff,
	tracking map[int]SlurTieGroup,
	groups *[]SlurTieGroup,
) {
	if _, ok := tracking[sid]; !ok {
		tracking[sid] = newSlurTieGroup(yPos, staffLine)
	}

	temp := tracking[sid]
	pos := entity.NewCoordinate(float64(note.PositionX), yPos)

	switch slur.Type {
	case musicxml.NoteSlurTypeStop, musicxml.NoteSlurTypeHop:
		temp.End = pos
	case musicxml.NoteSlurTypeStart:
		temp.Start = pos
		temp.Slur = &slur
	}

	temp.NoteMember = append(temp.NoteMember, note.UUID)
	temp.AccumulativeDirection += direction
	if yPos != 0 {
		temp.MaxY = math.Max(temp.MaxY, yPos)
		temp.MinY = math.Min(temp.MinY, yPos)
	}

	tracking[sid] = temp

	if slur.Type == musicxml.NoteSlurTypeStop || slur.Type == musicxml.NoteSlurTypeHop {
		*groups = append(*groups, temp)
		delete(tracking, sid)
	}

	if slur.Type == musicxml.NoteSlurTypeHop {
		next := newSlurTieGroup(yPos, staffLine)
		next.NoteMember = []string{note.UUID}
		next.Start = entity.NewCoordinate(float64(note.PositionX), yPos)
		next.AccumulativeDirection = temp.AccumulativeDirection
		next.Slur = &slur
		tracking[sid] = next
	}
}

func GetYPosGroup(group []entity.CoordinateWithNoteLength, staffLine lines.LineStaff) int {
	staffMiddleLine := staffLine.GetMiddleLine()

	slices.SortFunc(group, func(a, b entity.CoordinateWithNoteLength) int {
		return cmp.Compare(math.Abs(a.Y-float64(staffMiddleLine)), math.Abs(b.Y-float64(staffMiddleLine)))
	})
	farthest := group[len(group)-1]

	return staffLine.GetStemDirectionCompare(farthest.Y)
}

func GetBezierPoints(st SlurTieGroup, lineStaff lines.LineStaff, direction int) (entity.Coordinate, entity.Coordinate, entity.Coordinate) {
	tiesEndOffset := map[int][2]entity.Coordinate{
		-1: {entity.NewCoordinate(6, -6), entity.NewCoordinate(6, -6)},
		1:  {entity.NewCoordinate(6, 6), entity.NewCoordinate(3, 6)},
	}
	getYVal := map[int]float64{-1: st.MinY, 1: st.MaxY}

	// TODO: handle nested slur inside slur. st.Slur.Number > 1
	start := st.Start
	end := st.End

	placeholderY := float64(lineStaff.GetMiddleLine())
	if start.Y != 0 {
		placeholderY = start.Y
	} else if end.Y != 0 {
		placeholderY = end.Y
	}

	// no start
	if start.X == 0 || start.Y == 0 {
		start = entity.NewCoordinate(float64(lineStaff.GetLeftIndent())-(20), placeholderY)
	}

	// no end
	if end.X == 0 || end.Y == 0 {
		end = entity.NewCoordinate(float64(lineStaff.GetMarginRight()), placeholderY)
	}

	pull := entity.NewCoordinate(
		start.X+((end.X-start.X)/2),
		getYVal[direction])

	if pull.Y <= 0 {
		minMax := map[int]func(float64, float64) float64{-1: math.Min, 1: math.Max}
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

	resultStart := entity.NewCoordinate(math.Round(start.X)+tiesEndOffset[direction][0].X, math.Round(start.Y)+tiesEndOffset[direction][0].Y)
	resultPull := entity.NewCoordinate(math.Round(pull.X)+3+(tiesEndOffset[direction][0].X/2), math.Ceil(pull.Y+pullY))
	resultEnd := entity.NewCoordinate(math.Round(end.X)+tiesEndOffset[direction][1].X, math.Round(end.Y)+tiesEndOffset[direction][1].Y)

	return resultStart, resultPull, resultEnd

}

func RenderSlurTies(canv canvas.Canvas, lineStaff lines.LineStaff, groupBeam [][]entity.CoordinateWithNoteLength, slurties []SlurTieGroup) [2]entity.Coordinate {
	canv.Group(`class="slurties"`)

	pulls := []entity.Coordinate{}

	tiesNotes := map[string]bool{}

	sort.Slice(slurties, func(i, j int) bool {
		return slurties[i].Ties != nil && slurties[j].Ties == nil
	})

	for _, st := range slurties {

		group := []entity.CoordinateWithNoteLength{}

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

			// TODO: remove duplicate. that overalap between the group and the slur ties
			direction = GetYPosGroup(group, lineStaff) //+ st.AccumulativeDirection
			if direction < 0 {
				direction = -1
			} else {
				direction = 1
			}
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

		start, pull, end := GetBezierPoints(st, lineStaff, direction)

		// handle ties inside slur
		if !isTies && tiesNotes[st.NoteMember[0]] {
			pull.Y += float64(direction) * 10
			start.Y += float64(direction) * 3
		}

		lineType := "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1.5"
		if sluLineType == musicxml.NoteSlurLineTypeDashed {
			// Calculate approximate curve length
			curveLen := utils.QuadBezierLength(start, pull, end, 30)

			dash := 3.5
			patternCount := math.Floor(curveLen / (dash * 2))
			gap := (curveLen / patternCount) - dash

			lineType += fmt.Sprintf(";stroke-dasharray:%.1f %.1f;", dash, gap)
			lineType += "stroke-dashoffset:" + fmt.Sprintf("%f", dash/2) + ";"

			canv.Qbez(int(start.X), int(start.Y), int(pull.X), int(pull.Y), int(end.X), int(end.Y),
				lineType, fmt.Sprintf(`direction="%d"`, direction),
			)
		} else {
			w := canv.Writer()
			fmt.Fprintf(w, `<path d="%s" style="stroke-width: 0.4;stroke: #000000;" direction="%d"></path>`,
				GenerateTaperedSlur(start.X, start.Y, pull.X, pull.Y, end.X, end.Y),
				direction)
		}

		pulls = append(pulls, pull)
	}
	canv.Gend()

	sort.Slice(pulls, func(i, j int) bool {
		return pulls[i].Y < pulls[j].Y
	})
	return [2]entity.Coordinate{pulls[0], pulls[len(pulls)-1]}

}

func GenerateTaperedSlur(x1, y1, cx, cy, x2, y2 float64) string {
	thickness := 4.0 // The look thickness

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
