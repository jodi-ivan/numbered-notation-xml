package gregorian

import (
	"cmp"
	"context"
	"fmt"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func getAccidental(key keysig.Key, note *entity.NoteRenderer, accidentalsBefore map[int]map[string]musicxml.NoteAccidental) string {
	if note.AbsoluteAccidental == "" {
		return ""
	}
	hexMap := accidentalHex
	if accidentalsBefore[note.MeasureNumber][note.GetNonAccidentalAbsoluteNote()] == note.AbsoluteAccidental {
		hexMap = accidentalHexWithParentheses
	}

	return hexMap[note.AbsoluteAccidental]
}

func renderBean(canv canvas.Canvas, pos entity.Coordinate, noteType musicxml.NoteLength, note string, accidental string, octave int, attrs ...string) {
	attrs = append(attrs, fmt.Sprintf(`pitch="%s"`, note), fmt.Sprintf(`octave="%d"`, octave))
	if accidental != "" {
		canv.TextUnescaped(pos.X-8, pos.Y, accidental, `style="font-size:0.8em"`)
	}
	canv.TextUnescaped(pos.X, pos.Y,
		beanNoteHex[noteType],
		attrs...)

}
func RenderNote(ctx context.Context, canv canvas.Canvas, lines [5]int, groupBeam [][]CoordinateWithNoteLength, slursties []SlurTieGroup, notePos int, notes []*entity.NoteRenderer, timeSignature timesig.TimeSignature, keySignature keysig.KeySignature) (VMargin, [][]CoordinateWithNoteLength, []SlurTieGroup) {
	initialY := lines[0]
	note := notes[notePos]
	// maxY := lines[4]
	margin := VMargin{
		Top:    entity.NewCoordinate(0, float64(lines[0])),
		Bottom: entity.NewCoordinate(0, float64(lines[4])),
	}

	pairs := []SlurTieGroup{}
	yPos := GetYpos(lines, STAFF_SPACE_WIDTH, note.AbsoluteOctave, rune(note.AbsoluteNote[0]))
	margin.Set(entity.NewCoordinate(float64(note.PositionX), yPos))

	ts := timeSignature.GetTimesignatureOnMeasure(ctx, note.MeasureNumber)

	nonDottedValue := timeSignature.GetNoteLength(ctx, note.MeasureNumber, musicxml.Note{Type: note.NoteLength})
	dottedValue := timeSignature.GetNoteLength(ctx, note.MeasureNumber, musicxml.Note{Type: note.NoteLength, Dot: []*musicxml.Dot{{}}})

	merged := note.NoteLength == musicxml.NoteLengthEighth && notePos < len(notes)-1 && notes[notePos+1].IsDotted && notes[notePos+1].NoteLength == note.NoteLength
	dottedHalf := notePos < len(notes)-1 && notes[notePos+1].IsDotted && len(notes[notePos+1].Beam) >= 1
	dottedBeat := note.NoteValue == dottedValue
	quarterNoteInCompound := ts.IsCompoundTime() && note.NoteLength == musicxml.NoteLengthQuarter && dottedHalf

	merged = merged || quarterNoteInCompound

	beamType := note.NoteLength

	hasRemainingNote := note.NoteValue != dottedValue && note.NoteValue != nonDottedValue
	sameNextTwoDottedReplaced := notePos+2 < len(notes) && (notes[notePos+2].IsDotted || notes[notePos+2].AbsoluteNote == note.AbsoluteNote)
	sameNextDotted := notePos+1 < len(notes) && (notes[notePos+1].IsDotted || notes[notePos+1].NoteLength == note.NoteLength)
	sameNexTwoDotted := notePos+2 < len(notes) && (notes[notePos+2].IsDotted || notes[notePos+2].NoteLength == note.NoteLength)

	hasTiedNotes := sameNextTwoDottedReplaced || sameNextDotted || sameNexTwoDotted
	accidental := getAccidental(keySignature.GetKeyOnMeasure(ctx, note.MeasureNumber), note, nil)

	direction := cmp.Compare(yPos, float64(lines[2]))

	accumulative := 0
	memberGroup := 0
	if acc, ok := getDirectionAccumulative(slursties, note.UUID); ok {
		accumulative += acc
		memberGroup++
	}

	// if the current does not belong in any group, dont apply group consensus
	if memberGroup > 0 {
		if accumulative >= 0 {
			direction = 1
		} else {
			direction = -1
		}
	}

	note.StemDirection = direction

	if !merged && hasRemainingNote && hasTiedNotes {
		// since this is inserted additional notes, so the grouping as ties is guaranteed
		// just use the the existing group consensus.
		// this only enforce the stem, needs mechanism for the rendering slurties direction
		if accumulative >= 0 {
			direction = 1
		} else {
			direction = -1
		}

		note.StemDirection = direction

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
				canv.TextUnescaped(float64(xPos+12), dotPos, "&#xF060;", `style="fill:#0000DD"`)
			}
		}

		renderBean(
			canv,
			entity.NewCoordinate(float64(xPos), yPos), noteType[remaining],
			note.AbsoluteNote, accidental, note.AbsoluteOctave,
			fmt.Sprintf(`value="%.f"`, note.NoteValue), `style="fill:#0000DD"`)

		info := renderStemAndBeamMap[direction](canv, lines, CoordinateWithNoteLength{
			Coordinate: entity.NewCoordinate(float64(xPos), yPos),
			NoteLength: note.NoteLength,
			Beam:       note.Beam,
			NoteID:     note.UUID,
		})

		margin.SetBottom(info.LowestYPosition)
		margin.SetTop(info.HighestYPosition)

		pair := SlurTieGroup{
			AccumulativeDirection: direction,
			NoteMember: []string{
				note.UUID,
				notes[notePos+nextNotePos].UUID,
			},
			Start: entity.NewCoordinate(float64(note.PositionX), yPos),
			End:   entity.NewCoordinate(float64(notes[notePos+nextNotePos].PositionX), yPos),
			Ties: &entity.Slur{
				Number: 1,
			},
		}

		pairs = append(pairs, pair)

		if remaining < 1 {
			groupBeam = append(groupBeam, []CoordinateWithNoteLength{
				{
					Coordinate: entity.NewCoordinate(float64(xPos), yPos),
					NoteLength: noteType[remaining],
					NoteID:     notes[notePos+nextNotePos].UUID,
				},
			})
		}
	}

	RenderLedgerLine(canv, entity.NewCoordinate(float64(note.PositionX), yPos), lines)
	xPos := float64(note.PositionX)

	marker := `style="fill:#DD0000"`
	if note.NoteValue == nonDottedValue || note.NoteValue == dottedValue {
		marker = ""
	}
	renderBean(canv,
		entity.NewCoordinate(xPos, yPos),
		beamType, note.AbsoluteNote, accidental, note.AbsoluteOctave,
		fmt.Sprintf(`value="%.f"`, note.NoteValue), fmt.Sprintf(`uuid="%s"`, note.UUID), marker)

	if (dottedHalf || dottedBeat) && (!merged || (merged && quarterNoteInCompound && note.NoteValue-1 == 2)) {
		dotPos := yPos
		if (int(yPos)-initialY)%STAFF_SPACE_WIDTH == 0 {
			dotPos -= 4
		}
		canv.TextUnescaped(xPos+12, dotPos, "&#xF060;")
	}

	if note.NoteLength == musicxml.NoteLengthWhole {

		return margin, groupBeam, pairs
	}

	nextNoteIsDotted := notePos < len(notes)-1 && notes[notePos+1].IsDotted
	nextNoteIsSameNoteLength := notePos < len(notes)-1 && notes[notePos+1].NoteLength == note.NoteLength

	mergeNote := note.NoteLength == musicxml.NoteLengthEighth && nextNoteIsDotted && nextNoteIsSameNoteLength
	mergeNote = mergeNote || quarterNoteInCompound
	if len(note.Beam) == 0 || mergeNote {
		canv.Group(`group="false"`, fmt.Sprintf(`follow-consensus="%v"`, memberGroup > 0), fmt.Sprintf(`direction="%d"`, accumulative))

		stemInfo := renderStemAndBeamMap[direction](canv, lines, CoordinateWithNoteLength{
			Coordinate: entity.NewCoordinate(xPos, yPos),
			NoteLength: note.NoteLength,
			Beam:       note.Beam,
		})
		canv.Gend()

		// since this is only vertical stem line. only bother to move if it really disruptive like will invade the numbered note space. except it has beam.
		if margin.Bottom.Y+(STAFF_SPACE_WIDTH) < stemInfo.LowestYPosition.Y || (margin.Bottom.Y < stemInfo.LowestYPosition.Y && len(note.Beam) > 0) {
			margin.SetBottom(entity.NewCoordinate(float64(note.PositionX), stemInfo.LowestYPosition.Y))
		}

		margin.SetTop(stemInfo.HighestYPosition)

		if mergeNote {
			groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
		}

		return margin, groupBeam, pairs
	}

	groupBeam[len(groupBeam)-1] = append(groupBeam[len(groupBeam)-1], CoordinateWithNoteLength{
		Coordinate: entity.NewCoordinate(xPos, yPos),
		NoteLength: note.NoteLength,
		Beam:       note.Beam,
		NoteID:     note.UUID,
	})
	if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
		groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
	}
	return margin, groupBeam, pairs

}
