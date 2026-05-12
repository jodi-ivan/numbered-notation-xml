package gregorian

import (
	"context"
	"fmt"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	sline "github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
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

func RenderStaffLine(ctx context.Context, staffPos, y int, canv canvas.Canvas, notes []*entity.NoteRenderer, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature) VMargin {
	canv.Group(`class="gregorian"`, "style='font-family:mozart11'")

	lineStaff := sline.NewLineStaff(timeSignature, keySignature)
	lineStaff.Render(canv, y, notes[0].MeasureNumber, staffPos == 0)
	margin := VMargin{
		Top:           entity.NewCoordinate(0, float64(lineStaff.GetTopLine())),
		Bottom:        entity.NewCoordinate(0, float64(lineStaff.GetBottomLine())),
		DefaultTop:    lineStaff.GetTopLine(),
		DefaultBottom: lineStaff.GetBottomLine(),
	}

	groupBeam := [][]CoordinateWithNoteLength{{}}

	canv.Group(`class="notes"`, `style="font-size:2em"`)
	currentMeasure := 0

	groupBeamSlurTies := GetGroupSlueTies(notes, lineStaff)

	for i, note := range notes {

		if currentMeasure != note.MeasureNumber {
			currentMeasure = note.MeasureNumber
			if i != 0 {
				canv.Gend()
			}
			canv.Group(`class="measure"`, fmt.Sprintf(`number="%d"`, currentMeasure))

		}
		if note.IsAdditional {
			continue
		}

		if breathpause.IsBreathMark(note) {

			xPos := float64(note.PositionX)
			if note.PositionX-notes[i-1].PositionX <= numbered.MIN_DISTANCE_BREATH {
				xPos += (numbered.AVERAGE_CHARACTER_WIDTH + constant.LOWERCASE_LENGTH) / 3
			}

			canv.TextUnescaped(xPos, float64(lineStaff.GetTopLine())-STAFF_SPACE_WIDTH, "&#xF0E2;", `style="font-size:1.3em"`)
			if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
				groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
			}
			continue
		}
		if note.IsRest {
			canv.TextUnescaped(float64(note.PositionX), float64(lineStaff.GetMiddleLine()), restHex[note.NoteLength])
			if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
				groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
			}
			continue
		}

		if note.Barline != nil {
			barlinePos := entity.NewCoordinate(float64(note.PositionX), float64(lineStaff.GetBottomLine()))
			barline.RenderGregorian(canv, note.Barline, i == len(notes)-1, lineStaff, barlinePos)
			continue
		}

		if note.AbsoluteNote == "" {
			continue
		}

		var noteMargin VMargin
		pairs := []SlurTieGroup{}
		noteMargin, groupBeam, pairs = RenderNote(ctx, canv, lineStaff, groupBeam, groupBeamSlurTies, i, notes, timeSignature, keySignature)
		margin.Merge(noteMargin)

		groupBeamSlurTies = append(groupBeamSlurTies, pairs...)

	}
	canv.Gend()

	directions := map[int][]CoordinateWithNoteLength{
		1: {}, -1: {},
	}

	for i, gr := range groupBeam {
		if len(gr) == 0 {
			continue
		}
		gMargin, direction := RenderGroupBeam(canv, gr, lineStaff, groupBeamSlurTies)
		margin.Merge(gMargin)
		directions[direction] = append(directions[i], gr...)
	}

	for dir, locs := range directions {
		for _, loc := range locs {
			for _, note := range notes {
				if note.UUID == loc.NoteID {
					note.StemDirection = dir
				}
			}
		}

	}

	canv.Gend()

	st := RenderSlurTies(canv, lineStaff, groupBeam, groupBeamSlurTies)
	margin.Merge(st)

	canv.Gend()

	for _, note := range notes {
		if note.Note == 0 {
			continue
		}
		beanPos := lineStaff.GetYPos(rune(note.AbsoluteNote[0]), note.AbsoluteOctave)
		maxY := beanPos

		if note.StemDirection == 1 {
			maxY = beanPos + math.Floor(float64(note.StemDirection)*(2.5*sline.STAFF_SPACE_WIDTH))
		}

		if maxY < float64(lineStaff.GetTopLine()) {
			note.MarginTopFromStaff = lineStaff.GetTopLine() - int(maxY)
		}

		// REFACTOR THIS
		if note.Fermata != nil {
			breathpause.RenderFermata(ctx,
				canv, note.Fermata,
				entity.NewCoordinate(float64(note.PositionX), float64(y+10-note.MarginTopFromStaff)))

		}

	}

	return margin
}
