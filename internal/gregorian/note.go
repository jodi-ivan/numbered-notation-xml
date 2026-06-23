package gregorian

import (
	"context"
	"fmt"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/rhythm"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
	stfline "github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
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

func renderBean(canv canvas.Canvas, pos entity.Coordinate, noteType musicxml.NoteLength, note string, accidental string, octave, topLineStaff int, dotted bool) {
	if accidental != "" {
		canv.TextUnescaped(pos.X-8, pos.Y, accidental, `style="font-size:25.6px"`)
	}
	canv.TextUnescaped(pos.X, pos.Y, beanNoteHex[noteType])

	if dotted {
		dotPos := pos.Y
		if (int(dotPos)-topLineStaff)%STAFF_SPACE_WIDTH == 0 {
			dotPos -= 4
		}
		canv.TextUnescaped(pos.X+12, dotPos, "&#xF060;")
	}

}

func assignStemDirection(directions map[int][]entity.CoordinateWithNoteLength, notes []*entity.NoteRenderer) {
	notesMap := map[string]*entity.NoteRenderer{}

	for _, note := range notes {
		notesMap[note.UUID] = note
	}

	// assign direction
	for dir, locs := range directions {
		for _, loc := range locs {
			notesMap[loc.NoteID].StemDirection = dir
		}
	}
}

// calculate margin top
func setMarginTop(notes []*entity.NoteRenderer, lineStaff lines.LineStaff) {

	for _, note := range notes {
		if note.Note == 0 {
			continue
		}
		beanPos := lineStaff.GetYPos(rune(note.AbsoluteNote[0]), note.AbsoluteOctave)
		maxY := beanPos

		if note.StemDirection == 1 {
			maxY = beanPos + math.Floor(float64(note.StemDirection)*(2.5*lines.STAFF_SPACE_WIDTH))
		}

		if maxY < float64(lineStaff.GetTopLine()) {
			note.MarginTopFromStaff = lineStaff.GetTopLine() - int(maxY)
		}

	}
}

// getAdditionalNotes when a whole note length cant be represented as one note
// we need to add additional notes and register both of them in ties.
func getAdditionalNotes(ctx context.Context, ts timesig.TimeSignature, notes []*entity.NoteRenderer, staffLines lines.LineStaff, notePos, accumulativeDirection int, groupBeam *[][]entity.CoordinateWithNoteLength) *entity.NoteRenderer {
	note := notes[notePos]

	yPos := staffLines.GetYPos(rune(note.AbsoluteNote[0]), note.AbsoluteOctave)

	// check it can represent by a note
	nonDottedValue := ts.GetNoteLength(ctx, note.MeasureNumber, musicxml.Note{Type: note.NoteLength})
	dottedValue := ts.GetNoteLength(ctx, note.MeasureNumber, musicxml.Note{Type: note.NoteLength, Dot: []*musicxml.Dot{{}}})
	hasRemainingNote := note.NoteValue != dottedValue && note.NoteValue != nonDottedValue
	if !hasRemainingNote {
		return nil
	}

	// check the dots the next and next two notes.
	sameNextTwoDottedReplaced := notePos+2 < len(notes) && (notes[notePos+2].AbsoluteNote == note.AbsoluteNote) && notes[notePos+2].Tie != nil
	sameNextDotted := notePos+1 < len(notes) && (notes[notePos+1].IsDotted || notes[notePos+1].NoteLength == note.NoteLength)
	sameNexTwoDotted := notePos+2 < len(notes) && (notes[notePos+2].IsDotted || notes[notePos+2].NoteLength == note.NoteLength)

	hasTiedNotes := sameNextTwoDottedReplaced || sameNextDotted || sameNexTwoDotted
	if !hasTiedNotes {
		return nil
	}

	// merged that dotted in the next note numbered can be represented as one note in standard notation
	// half notes that represented as numbered and dotted. --> represented as a quarter note in standard notation
	merged := note.NoteLength == musicxml.NoteLengthEighth && notePos < len(notes)-1 && notes[notePos+1].IsDotted && notes[notePos+1].NoteLength == note.NoteLength

	currentKeysig := ts.GetTimesignatureOnMeasure(ctx, note.MeasureNumber)
	dottedHalf := notePos < len(notes)-1 && notes[notePos+1].IsDotted && len(notes[notePos+1].Beam) >= 1
	quarterNoteInCompound := currentKeysig.IsCompoundTime() && note.NoteLength == musicxml.NoteLengthQuarter && dottedHalf

	merged = merged || quarterNoteInCompound

	if merged {
		return nil
	}

	// since this is inserted additional notes, so the grouping as ties is guaranteed
	// just use the the existing group consensus.
	// this only enforce the stem, needs mechanism for the rendering slurties direction
	direction := -1
	if accumulativeDirection >= 0 {
		direction = 1
	}

	note.StemDirection = direction

	remaining := note.NoteValue - nonDottedValue
	nextNotePos := 2
	xPos := notes[notePos+2].PositionX

	if note.NoteLength == musicxml.NoteLengthEighth {
		xPos = notes[notePos+1].PositionX
		nextNotePos = 1
	}

	result := &entity.NoteRenderer{
		AbsoluteNote:       note.AbsoluteNote,
		AbsoluteAccidental: note.AbsoluteAccidental,
		AbsoluteOctave:     note.AbsoluteOctave,
		PositionX:          xPos,
		UUID:               notes[notePos+nextNotePos].UUID,
		NoteLength:         noteType[remaining],
	}

	if _, ok := noteType[remaining]; !ok {
		if remaining == 1.5 {
			remaining = 1
			result.IsDotted = true
		}
	}

	result.NoteLength = noteType[remaining]

	if remaining < 1 {
		*groupBeam = append(*groupBeam, []entity.CoordinateWithNoteLength{
			{
				Coordinate: entity.NewCoordinate(float64(xPos), yPos),
				NoteLength: noteType[remaining],
				NoteID:     notes[notePos+nextNotePos].UUID,
			},
		})
	}

	return result
}

func isDottedNote(notes []*entity.NoteRenderer, notePos int, ts timesig.TimeSignature) bool {

	note := notes[notePos]
	dottedValue := ts.GetNoteLength(context.Background(), note.MeasureNumber, musicxml.Note{Type: note.NoteLength, Dot: []*musicxml.Dot{{}}})

	dottedHalf := notePos < len(notes)-1 && notes[notePos+1].IsDotted && len(notes[notePos+1].Beam) >= 1
	dottedBeat := note.NoteValue == dottedValue

	merged := note.NoteLength == musicxml.NoteLengthEighth && notePos < len(notes)-1 && notes[notePos+1].IsDotted && notes[notePos+1].NoteLength == note.NoteLength

	currentKeysig := ts.GetTimesignatureOnMeasure(context.Background(), note.MeasureNumber)
	quarterNoteInCompound := currentKeysig.IsCompoundTime() && note.NoteLength == musicxml.NoteLengthQuarter && dottedHalf

	merged = merged || quarterNoteInCompound

	return (dottedHalf || dottedBeat) && (!merged || (merged && quarterNoteInCompound && note.NoteValue-1 == 2))
}

func RenderNote(ctx context.Context, canv canvas.Canvas, staffLines stfline.LineStaff, groupBeam [][]entity.CoordinateWithNoteLength, slursties []rhythm.SlurTieGroup, notePos int, notes []*entity.NoteRenderer, timeSignature timesig.TimeSignature, keySignature keysig.KeySignature) (VMargin, [][]entity.CoordinateWithNoteLength, []rhythm.SlurTieGroup) {

	canv.Group(`class="note"`, `style="font-size:32px"`)
	defer func() {
		canv.Gend()
	}()
	lines := staffLines.GetLines()
	initialY := staffLines.GetTopLine()

	note := notes[notePos]

	margin := VMargin{
		Top:    entity.NewCoordinate(0, float64(initialY)),
		Bottom: entity.NewCoordinate(0, float64(staffLines.GetBottomLine())),
	}

	pairs := []rhythm.SlurTieGroup{}
	yPos := staffLines.GetYPos(rune(note.AbsoluteNote[0]), note.AbsoluteOctave)
	currentPos := entity.NewCoordinate(float64(note.PositionX), yPos)
	margin.Set(currentPos)

	accidental := getAccidental(keySignature.GetKeyOnMeasure(ctx, note.MeasureNumber), note, nil)

	beamType := note.NoteLength
	direction := staffLines.GetStemDirectionCompare(yPos)

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

	additionalNotes := getAdditionalNotes(ctx, timeSignature, notes, staffLines, notePos, accumulative, &groupBeam)
	if additionalNotes != nil {
		// since this is inserted additional notes, so the grouping as ties is guaranteed
		// just use the the existing group consensus.
		// this only enforce the stem, needs mechanism for the rendering slurties direction
		if accumulative >= 0 {
			direction = 1
		} else {
			direction = -1
		}
		addtnlPos := entity.NewCoordinate(float64(additionalNotes.PositionX), yPos)
		pair := rhythm.SlurTieGroup{
			AccumulativeDirection: direction,
			NoteMember: []string{
				note.UUID,
				additionalNotes.UUID,
			},
			Start: entity.NewCoordinate(float64(note.PositionX), yPos),
			End:   addtnlPos,
			Ties: &entity.Slur{
				Number: 1,
			},
		}

		pairs = append(pairs, pair)

		renderBean(canv,
			addtnlPos,
			additionalNotes.NoteLength, note.AbsoluteNote, accidental, note.AbsoluteOctave, staffLines.GetTopLine(), additionalNotes.IsDotted)

		renderStemAndBeamMap[direction](canv, lines, entity.CoordinateWithNoteLength{
			Coordinate: addtnlPos,
			NoteLength: note.NoteLength,
			Beam:       note.Beam,
			NoteID:     note.UUID,
			Tuplet:     note.Tuplet,
		})

		RenderLedgerLine(canv, addtnlPos, staffLines.GetTopLine(), staffLines.GetBottomLine())

	}

	xPos := float64(note.PositionX)

	isDotted := isDottedNote(notes, notePos, timeSignature)
	renderBean(canv,
		entity.NewCoordinate(xPos, yPos),
		beamType, note.AbsoluteNote, accidental, note.AbsoluteOctave, staffLines.GetTopLine(), isDotted)
	RenderLedgerLine(canv, currentPos, staffLines.GetTopLine(), staffLines.GetBottomLine())

	if note.NoteLength == musicxml.NoteLengthWhole {

		return margin, groupBeam, pairs
	}

	nextNoteIsDotted := notePos < len(notes)-1 && notes[notePos+1].IsDotted
	nextNoteIsSameNoteLength := notePos < len(notes)-1 && notes[notePos+1].NoteLength == note.NoteLength

	dottedHalf := notePos < len(notes)-1 && notes[notePos+1].IsDotted && len(notes[notePos+1].Beam) >= 1
	ts := timeSignature.GetTimesignatureOnMeasure(ctx, note.MeasureNumber)
	quarterNoteInCompound := ts.IsCompoundTime() && note.NoteLength == musicxml.NoteLengthQuarter && dottedHalf

	mergeNote := note.NoteLength == musicxml.NoteLengthEighth && nextNoteIsDotted && nextNoteIsSameNoteLength
	mergeNote = mergeNote || quarterNoteInCompound
	if len(note.Beam) == 0 || mergeNote {
		canv.Group(fmt.Sprintf(`follow-consensus="%v"`, memberGroup > 0), fmt.Sprintf(`direction="%d"`, accumulative))

		stemInfo := renderStemAndBeamMap[direction](canv, lines, entity.CoordinateWithNoteLength{
			Coordinate: entity.NewCoordinate(xPos, yPos),
			NoteLength: note.NoteLength,
			Beam:       note.Beam,
			Tuplet:     note.Tuplet,
		})
		canv.Gend()

		// since this is only vertical stem line. only bother to move if it really disruptive like will invade the numbered note space. except it has beam.
		if margin.Bottom.Y+(STAFF_SPACE_WIDTH) < stemInfo.LowestYPosition.Y || (margin.Bottom.Y < stemInfo.LowestYPosition.Y && len(note.Beam) > 0) {
			margin.SetBottom(entity.NewCoordinate(float64(note.PositionX), stemInfo.LowestYPosition.Y))
		}

		margin.SetTop(stemInfo.HighestYPosition)

		if mergeNote {
			groupBeam = append(groupBeam, []entity.CoordinateWithNoteLength{})
		}

		return margin, groupBeam, pairs
	}

	groupBeam[len(groupBeam)-1] = append(groupBeam[len(groupBeam)-1], entity.CoordinateWithNoteLength{
		Coordinate: entity.NewCoordinate(xPos, yPos),
		NoteLength: note.NoteLength,
		Beam:       note.Beam,
		NoteID:     note.UUID,
		Tuplet:     note.Tuplet,
	})
	if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
		groupBeam = append(groupBeam, []entity.CoordinateWithNoteLength{})
	}
	return margin, groupBeam, pairs

}
