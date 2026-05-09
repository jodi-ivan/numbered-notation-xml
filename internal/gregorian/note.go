package gregorian

import (
	"cmp"
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func RenderNote(ctx context.Context, canv canvas.Canvas, lines [5]int, groupBeam [][]CoordinateWithNoteLength, notePos int, notes []*entity.NoteRenderer, timeSignature timesig.TimeSignature) (int, [][]CoordinateWithNoteLength) {
	initialY := lines[0]
	note := notes[notePos]
	maxY := lines[4]
	yPos := GetYpos(lines, STAFF_SPACE_WIDTH, note.AbsoluteOctave, rune(note.AbsoluteNote[0]))
	if maxY < int(yPos) {
		maxY = int(yPos)
	}

	ts := timeSignature.GetTimesignatureOnMeasure(ctx, note.MeasureNumber)

	dottedValue := timeSignature.GetNoteLength(ctx, note.MeasureNumber, musicxml.Note{Type: note.NoteLength, Dot: []*musicxml.Dot{{}}})
	merged := notePos < len(notes)-1 && note.NoteLength == musicxml.NoteLengthEighth && notes[notePos+1].NoteLength == note.NoteLength
	dottedHalf := notePos < len(notes)-1 && notes[notePos+1].IsDotted && len(notes[notePos+1].Beam) >= 1
	dottedBeat := note.NoteValue == dottedValue
	quarterNoteInCompound := ts.IsCompoundTime() && note.NoteLength == musicxml.NoteLengthQuarter && dottedHalf

	merged = merged || quarterNoteInCompound

	beamType := note.NoteLength

	// currNoteLength := ts.GetNoteLength(ctx, musicxml.Note{Type: note.NoteLength})

	nonDottedValue := timeSignature.GetNoteLength(ctx, note.MeasureNumber, musicxml.Note{Type: note.NoteLength})
	if !merged && note.NoteValue != dottedValue && note.NoteValue != nonDottedValue && notePos+2 < len(notes) && (notes[notePos+2].IsDotted || notes[notePos+2].AbsoluteNote == note.AbsoluteNote) {
		remaining := note.NoteValue - nonDottedValue
		nextNotePos := 2
		xPos := notes[notePos+2].PositionX

		if note.NoteLength == musicxml.NoteLengthEighth {
			xPos = notes[notePos+1].PositionX
			nextNotePos = 1
		}
		noteType := map[float64]musicxml.NoteLength{
			0.25: musicxml.NoteLength16th,
			0.5:  musicxml.NoteLengthEighth,
			1:    musicxml.NoteLengthQuarter,
			2:    musicxml.NoteLengthHalf,
		}

		if _, ok := noteType[remaining]; !ok {
			if remaining == 1.5 {
				remaining = 1
				dotPos := yPos
				if (int(yPos)-initialY)%STAFF_SPACE_WIDTH == 0 {
					dotPos -= 4
				}
				canv.TextUnescaped(float64(xPos+12), dotPos, "&#xF060;", `style="fill:#DD0000"`)
			}
		}

		canv.TextUnescaped(float64(xPos), yPos,
			beanNoteHex[noteType[remaining]],
			fmt.Sprintf(`pitch="%s"`, note.AbsoluteNote), fmt.Sprintf(`octave="%d"`, note.AbsoluteOctave), `style="fill:#DD0000"`)

		renderMap[cmp.Compare(yPos, float64(lines[2]))](canv, lines, CoordinateWithNoteLength{
			Coordinate: entity.NewCoordinate(float64(xPos), yPos),
			NoteLength: note.NoteLength,
			Beam:       note.Beam,
		})

		if note.Tie == nil {
			note.Tie = &entity.Slur{
				Number:        1,
				Type:          musicxml.NoteSlurTypeStart,
				GregorianOnly: true,
			}
		}

		if notes[notePos+nextNotePos].Tie != nil {
			notes[notePos+nextNotePos].Tie = &entity.Slur{
				Number:        1,
				Type:          musicxml.NoteSlurTypeStart,
				GregorianOnly: true,
			}
		}

		if remaining < 1 {
			groupBeam = append(groupBeam, []CoordinateWithNoteLength{
				{
					Coordinate: entity.NewCoordinate(float64(xPos), yPos),
					NoteLength: noteType[remaining],
				},
			})
		}
	}

	RenderLedgerLine(canv, entity.NewCoordinate(float64(note.PositionX), yPos), lines)

	canv.TextUnescaped(float64(note.PositionX), yPos,
		beanNoteHex[beamType],
		fmt.Sprintf(`pitch="%s"`, note.AbsoluteNote), fmt.Sprintf(`octave="%d"`, note.AbsoluteOctave))

	if (dottedHalf || dottedBeat) && (!merged || (merged && quarterNoteInCompound && note.NoteValue-1 == 2)) {
		dotPos := yPos
		if (int(yPos)-initialY)%STAFF_SPACE_WIDTH == 0 {
			dotPos -= 4
		}
		canv.TextUnescaped(float64(note.PositionX+12), dotPos, "&#xF060;")
	}

	if note.NoteLength == musicxml.NoteLengthWhole {

		return maxY, groupBeam
	}

	nextNoteIsDotted := notePos < len(notes)-1 && notes[notePos+1].IsDotted
	nextNoteIsSameNoteLength := notePos < len(notes)-1 && notes[notePos+1].NoteLength == note.NoteLength

	mergeNote := note.NoteLength == musicxml.NoteLengthEighth && nextNoteIsDotted && nextNoteIsSameNoteLength
	mergeNote = mergeNote || quarterNoteInCompound
	if len(note.Beam) == 0 || mergeNote {
		stemInfo := renderMap[cmp.Compare(yPos, float64(lines[2]))](canv, lines, CoordinateWithNoteLength{
			Coordinate: entity.NewCoordinate(float64(note.PositionX), yPos),
			NoteLength: note.NoteLength,
			Beam:       note.Beam,
		})

		// since this is only vertical stem line. only bother to move if it really disruptive like will invade the numbered note space. except it has beam.
		if maxY+(STAFF_SPACE_WIDTH*2) < int(stemInfo.LowestYPosition) || (maxY < int(stemInfo.LowestYPosition) && len(note.Beam) > 0) {
			maxY = int(stemInfo.LowestYPosition)
		}

		if mergeNote {
			groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
		}

		return maxY, groupBeam
	}

	groupBeam[len(groupBeam)-1] = append(groupBeam[len(groupBeam)-1], CoordinateWithNoteLength{
		Coordinate: entity.NewCoordinate(float64(note.PositionX), yPos),
		NoteLength: note.NoteLength,
		Beam:       note.Beam,
	})
	if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
		groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
	}
	return maxY, groupBeam

}
