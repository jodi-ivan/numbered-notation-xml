package gregorian

import (
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderStaffLine(ctx context.Context, staffPos, y int, canv canvas.Canvas, notes []*entity.NoteRenderer, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature) VMargin {

	lineStaff := lines.NewLineStaff(timeSignature, keySignature)
	lineStaff.Render(canv, y, notes[0].MeasureNumber, staffPos == 0)
	margin := VMargin{
		Top:           entity.NewCoordinate(0, float64(lineStaff.GetTopLine())),
		Bottom:        entity.NewCoordinate(0, float64(lineStaff.GetBottomLine())),
		DefaultTop:    lineStaff.GetTopLine(),
		DefaultBottom: lineStaff.GetBottomLine(),
	}

	groupBeam := [][]entity.CoordinateWithNoteLength{{}}

	canv.Group(`class="notes"`, `style="font-size:2em"`)
	currentMeasure := 0

	groupBeamSlurTies := rhythm.GetGroupSlueTies(notes, lineStaff)

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

			posX := note.PositionX
			prevNotePosX := notes[i-1].PositionX

			canv.TextUnescaped(
				breathpause.AdjustPosition(posX, prevNotePosX),
				float64(lineStaff.GetTopLine())-STAFF_SPACE_WIDTH,
				"&#xF0E2;", `style="font-size:1.3em"`)

			if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
				groupBeam = append(groupBeam, []entity.CoordinateWithNoteLength{})
			}
			continue
		}
		if note.IsRest {
			canv.TextUnescaped(float64(note.PositionX), float64(lineStaff.GetMiddleLine()), restHex[note.NoteLength])
			if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
				groupBeam = append(groupBeam, []entity.CoordinateWithNoteLength{})
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
		pairs := []rhythm.SlurTieGroup{}
		noteMargin, groupBeam, pairs = RenderNote(ctx, canv, lineStaff, groupBeam, groupBeamSlurTies, i, notes, timeSignature, keySignature)
		margin.Merge(noteMargin)

		groupBeamSlurTies = append(groupBeamSlurTies, pairs...)

	}
	canv.Gend()

	directions, gMargin := RenderGroupBeam(canv, lineStaff, groupBeam, groupBeamSlurTies)
	margin.Merge(gMargin)
	canv.Gend()

	assignStemDirection(directions, notes)

	if len(groupBeamSlurTies) > 0 {
		st := rhythm.RenderSlurTies(canv, lineStaff, groupBeam, groupBeamSlurTies)
		margin.Set(st[:]...)
	}
	setMarginTop(notes, lineStaff)

	// Fermata
	for _, note := range notes {

		if note.Fermata != nil {
			breathpause.RenderFermata(ctx,
				canv, note.Fermata,
				entity.NewCoordinate(float64(note.PositionX), float64(y+10-note.MarginTopFromStaff)))

		}

	}

	return margin
}
