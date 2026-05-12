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
	groupBeamSlurTies := []SlurTieGroup{}

	var tiesTracking *SlurTieGroup
	slurTracking := map[int]SlurTieGroup{}

	for _, note := range notes {
		yPos := 0.0
		direction := 0
		if note.AbsoluteNote != "" && note.AbsoluteOctave > 0 {
			yPos = staffLine.GetYPos(rune(note.AbsoluteNote[0]), note.AbsoluteOctave)
			direction = staffLine.GetStemDirection(rune(note.AbsoluteNote[0]), note.AbsoluteOctave)
		}

		if note.Tie != nil {
			if tiesTracking == nil && note.Tie.Type == musicxml.NoteSlurTypeStart && !note.Tie.NumberedOnly {
				tiesTracking = &SlurTieGroup{
					MaxY: yPos, MinY: yPos,
					Ties:       note.Tie,
					NoteMember: []string{note.UUID},
					Start:      entity.NewCoordinate(float64(note.PositionX), yPos),
				}
				tiesTracking.NoteMember = append(tiesTracking.NoteMember, note.UUID)
				tiesTracking.AccumulativeDirection += direction
			}

			if tiesTracking != nil && note.Tie.Type == musicxml.NoteSlurTypeStop {
				tiesTracking.AccumulativeDirection += direction

				tiesTracking.End = entity.NewCoordinate(float64(note.PositionX), yPos)
				tiesTracking.NoteMember = append(tiesTracking.NoteMember, note.UUID)

				groupBeamSlurTies = append(groupBeamSlurTies, *tiesTracking)
				tiesTracking = nil
			}
		}

		for sid, slur := range note.Slur {
			_, ok := slurTracking[sid]

			if !ok {
				if yPos == 0 {
					slurTracking[sid] = SlurTieGroup{
						MaxY: float64(staffLine.GetTopLine()),
						MinY: float64(staffLine.GetBottomLine()),
					}
				} else {
					slurTracking[sid] = SlurTieGroup{
						MaxY: yPos,
						MinY: yPos,
					}
				}
			}

			temp := slurTracking[sid]
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

			slurTracking[sid] = temp

			if slur.Type == musicxml.NoteSlurTypeStop || slur.Type == musicxml.NoteSlurTypeHop {
				groupBeamSlurTies = append(groupBeamSlurTies, temp)
				delete(slurTracking, sid)
			}

			if slur.Type == musicxml.NoteSlurTypeHop {
				temp := SlurTieGroup{
					NoteMember:            []string{note.UUID},
					Start:                 entity.NewCoordinate(float64(note.PositionX), yPos),
					AccumulativeDirection: temp.AccumulativeDirection,
					Slur:                  &slur,
				}

				if yPos == 0 {
					temp.MaxY = float64(staffLine.GetTopLine())
					temp.MinY = float64(staffLine.GetBottomLine())
				} else {
					temp.MaxY = yPos
					temp.MinY = yPos
				}
				slurTracking[sid] = temp

			}
		}

		if len(note.Slur) > 0 {
			continue
		}

		for i, v := range slurTracking {
			v.AccumulativeDirection += direction
			if yPos != 0 {
				v.MaxY = math.Max(v.MaxY, yPos)
				v.MinY = math.Min(v.MinY, yPos)
			}
			v.NoteMember = append(v.NoteMember, note.UUID)
			slurTracking[i] = v
		}

		if tiesTracking != nil && yPos != 0 {
			tiesTracking.MaxY = math.Max(tiesTracking.MaxY, yPos)
			tiesTracking.MinY = math.Min(tiesTracking.MinY, yPos)
		}

	}

	if tiesTracking != nil {
		groupBeamSlurTies = append(groupBeamSlurTies, *tiesTracking)

	}
	for _, v := range slurTracking {
		groupBeamSlurTies = append(groupBeamSlurTies, v)
	}
	return groupBeamSlurTies
}

func GetYPosGroup(group []entity.CoordinateWithNoteLength, staffLine lines.LineStaff) int {
	staffMiddleLine := staffLine.GetMiddleLine()

	slices.SortFunc(group, func(a, b entity.CoordinateWithNoteLength) int {
		return cmp.Compare(math.Abs(a.Y-float64(staffMiddleLine)), math.Abs(b.Y-float64(staffMiddleLine)))
	})
	farthest := group[len(group)-1]

	return staffLine.GetStemDirectionCompare(farthest.Y)
}

func RenderSlurTies(canv canvas.Canvas, lineStaff lines.LineStaff, groupBeam [][]entity.CoordinateWithNoteLength, slurties []SlurTieGroup) [2]entity.Coordinate {
	canv.Group(`class="slurties"`)

	pulls := []entity.Coordinate{}

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
			direction = GetYPosGroup(group, lineStaff)
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

		placeholderY := float64(lineStaff.GetMiddleLine())
		if start.Y != 0 {
			placeholderY = start.Y
		} else if end.Y != 0 {
			placeholderY = end.Y
		}

		// no start
		if start.X == 0 || start.Y == 0 {
			start = entity.NewCoordinate(float64(lineStaff.GetLeftIndent())-(10*2), placeholderY)
		}

		// no end
		if end.X == 0 || end.Y == 0 {
			end = entity.NewCoordinate(float64(lineStaff.GetMarginRight()), placeholderY)
		}

		getYVal := map[int]float64{
			-1: st.MinY,
			1:  st.MaxY,
		}

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

		pulls = append(pulls, entity.NewCoordinate(math.Round(pull.X)+3+(tiesEndOffset[direction][0].X/2), pull.Y+pullY))
	}
	canv.Gend()

	sort.Slice(pulls, func(i, j int) bool {
		return pulls[i].Y < pulls[j].Y
	})
	return [2]entity.Coordinate{pulls[0], pulls[len(pulls)-1]}

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
