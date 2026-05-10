package gregorian

import (
	"cmp"
	"context"
	"fmt"
	"unicode"

	"github.com/jodi-ivan/numbered-notation-xml/internal/barline"
	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/numbered"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

func GetYPosKeySig(lines [5]int, space float64, pitch string, isFlat bool) float64 {
	if isFlat {
		// Flat order: B E A D G C F
		pos := map[string]float64{
			"B": float64(lines[2]),
			"E": float64(lines[0]) + (space / 2),
			"A": float64(lines[2]) + (space / 2),
			"D": float64(lines[1]),
			"G": float64(lines[3]),
			"C": float64(lines[1]) + (space / 2),
			"F": float64(lines[3]) - (space / 2),
		}
		return pos[pitch]
	}
	// Sharp order: F C G D A E B
	pos := map[string]float64{
		"F": float64(lines[0]),
		"C": float64(lines[2]) - (space / 2),
		"G": float64(lines[1]) - (space / 2),
		"D": float64(lines[2]),
		"A": float64(lines[3]) + (space / 2),
		"E": float64(lines[0]) + (space / 2),
		"B": float64(lines[2]) - (space),
	}
	return pos[pitch]
}

func GetYpos(lines [5]int, space int, octave int, pitch rune) float64 {
	noteOrder := []rune{'C', 'D', 'E', 'F', 'G', 'A', 'B'}

	diatonicIndex := func(p rune, oct int) int {
		for i, n := range noteOrder {
			if n == unicode.ToUpper(p) {
				return oct*7 + i
			}
		}
		return -1
	}

	refIndex := diatonicIndex('F', 5) // lines[0] = F5
	noteIndex := diatonicIndex(pitch, octave)

	stepsBelow := refIndex - noteIndex

	return float64(lines[0]) + float64(stepsBelow)*(float64(space)/2)
}

func GetGroupSlueTies(notes []*entity.NoteRenderer, lines [5]int) []SlurTieGroup {
	groupBeamSlurTies := []SlurTieGroup{}

	var tiesTracking *SlurTieGroup
	slurTracking := map[int]SlurTieGroup{}

	for _, note := range notes {
		yPos := 0.0
		direction := 0
		if note.AbsoluteNote != "" && note.AbsoluteOctave > 0 {
			yPos = GetYpos(lines, STAFF_SPACE_WIDTH, note.AbsoluteOctave, rune(note.AbsoluteNote[0]))
			direction = cmp.Compare(int(yPos), lines[2])
			if direction == 0 {
				direction = -1
			}
		}

		if note.Tie != nil {
			if tiesTracking == nil && note.Tie.Type == musicxml.NoteSlurTypeStart {
				tiesTracking = &SlurTieGroup{
					Ties:       note.Tie,
					NoteMember: []string{},
				}
			}
			tiesTracking.NoteMember = append(tiesTracking.NoteMember, note.UUID)
			tiesTracking.AccumulativeDirection += direction

			if note.Tie.Type == musicxml.NoteSlurTypeStop {
				groupBeamSlurTies = append(groupBeamSlurTies, *tiesTracking)
				tiesTracking = nil
			}
		}

		for sid, slur := range note.Slur {
			_, ok := slurTracking[sid]

			if !ok {
				slurTracking[sid] = SlurTieGroup{}
			}

			temp := slurTracking[sid]
			pos := entity.NewCoordinate(float64(note.PositionX), yPos)

			switch slur.Type {
			case musicxml.NoteSlurTypeStop:
				temp.End = pos
			case musicxml.NoteSlurTypeStart:
				temp.Start = pos
			}

			temp.Slur = &slur
			temp.NoteMember = append(temp.NoteMember, note.UUID)
			temp.AccumulativeDirection += direction
			slurTracking[sid] = temp

			if slur.Type == musicxml.NoteSlurTypeStop || slur.Type == musicxml.NoteSlurTypeHop {
				groupBeamSlurTies = append(groupBeamSlurTies, temp)
				delete(slurTracking, sid)
			}

			if slur.Type == musicxml.NoteSlurTypeHop {
				temp := SlurTieGroup{
					NoteMember:            []string{note.UUID},
					Start:                 entity.NewCoordinate(float64(note.PositionX), yPos),
					AccumulativeDirection: direction,
				}
				slurTracking[sid] = temp

			}
		}

		if len(note.Slur) > 0 {
			continue
		}

		for i, v := range slurTracking {
			v.AccumulativeDirection += direction
			v.NoteMember = append(v.NoteMember, note.UUID)
			slurTracking[i] = v
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

func RenderStaffLine(ctx context.Context, staffPos, y int, canv canvas.Canvas, notes []*entity.NoteRenderer, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature) int {
	initialY := y
	lines := [5]int{}
	canv.Group(`class="gregorian"`, "style='font-family:mozart11'")
	x2 := constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH + 8
	canv.Group(`class="staff-line"`)
	for i := 0; i <= 4; i++ {
		lines[i] = y
		canv.Line(constant.LAYOUT_INDENT_LENGTH, y, x2, y, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.8")
		y += STAFF_SPACE_WIDTH
	}
	canv.Line(constant.LAYOUT_INDENT_LENGTH, initialY, constant.LAYOUT_INDENT_LENGTH, y-STAFF_SPACE_WIDTH, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1")
	canv.Gend()

	maxY := lines[4]

	groupBeam := [][]CoordinateWithNoteLength{{}}

	canv.Group(`class="notes"`, `style="font-size:2em"`)
	currentMeasure := 0

	groupBeamSlurTies := GetGroupSlueTies(notes, lines)

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

			canv.TextUnescaped(xPos, float64(lines[0])-STAFF_SPACE_WIDTH, "&#xF0E2;", `style="font-size:1.3em"`)
			if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
				groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
			}
			continue
		}
		if note.IsRest {
			canv.TextUnescaped(float64(note.PositionX), float64(lines[2]), restHex[note.NoteLength])
			if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
				groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
			}
			continue
		}

		if note.Barline != nil {
			barline.RenderGregorian(canv, note.Barline, i == len(notes)-1, lines, entity.NewCoordinate(float64(note.PositionX), float64(lines[4])))
			continue
		}

		if note.AbsoluteNote == "" {
			continue
		}

		var noteMaxYPos int
		pairs := []SlurTieGroup{}
		noteMaxYPos, groupBeam, pairs = RenderNote(ctx, canv, lines, groupBeam, groupBeamSlurTies, i, notes, timeSignature, keySignature)
		if maxY < noteMaxYPos {
			maxY = noteMaxYPos
		}

		groupBeamSlurTies = append(groupBeamSlurTies, pairs...)

	}
	canv.Gend()

	canv.Group()
	for _, gr := range groupBeam {
		if len(gr) == 0 {
			continue
		}
		groupMaxY := RenderGroupBeam(canv, gr, lines, groupBeamSlurTies)
		if maxY < groupMaxY {
			maxY = groupMaxY
		}
	}
	canv.Gend()

	canv.Gend()

	x := float64(constant.LAYOUT_INDENT_LENGTH)

	canv.Group(`class="staff-markings"`)
	// clef
	key := keySignature.GetKeyOnMeasure(ctx, notes[0].MeasureNumber)
	accidentalSet := key.GetAccidentals()

	canv.Group(`class="clef"`, `style="font-size:2em"`)
	canv.TextUnescaped(constant.LAYOUT_INDENT_LENGTH+5, float64(initialY+15), TREBLE_CLEF_HEX)
	canv.Gend()

	canv.Group(`class="keysig"`, `style="font-size:1.75em"`)
	offset := 0

	// key signature changes
	if key.Start && key.Prev != nil && notes[0].MeasureNumber != 1 {
		naturalSet := key.Prev.GetAccidentals()

		for x, acc := range naturalSet {
			accidental := accidentalHex[musicxml.NoteAccidentalNatural]
			width := ACCIDENTAL_KEY_SIGNATURE_WIDTH

			canv.TextUnescaped(float64(constant.LAYOUT_INDENT_LENGTH+CLEF_WIDTH)+float64(width*x),
				GetYPosKeySig(lines, STAFF_SPACE_WIDTH, acc, key.Prev.Fifth < 0),
				accidental)
		}

		offset = (len(naturalSet) * ACCIDENTAL_KEY_SIGNATURE_WIDTH) + PADDING_WIDTH
	}
	for x, acc := range accidentalSet {
		accidental := accidentalHex[musicxml.NoteAccidentalSharp]
		width := ACCIDENTAL_KEY_SIGNATURE_WIDTH
		if key.Fifth < 0 {
			accidental = accidentalHex[musicxml.NoteAccidentalFlat]
		}
		canv.TextUnescaped(float64(constant.LAYOUT_INDENT_LENGTH+CLEF_WIDTH+offset)+float64(width*x),
			GetYPosKeySig(lines, STAFF_SPACE_WIDTH, acc, key.Fifth < 0),
			accidental)
	}
	canv.Gend()

	x += CLEF_WIDTH + (float64(len(accidentalSet)) * ACCIDENTAL_KEY_SIGNATURE_WIDTH) + PADDING_WIDTH + float64(offset)

	if staffPos == 0 && len(timeSignature.Signatures) > 0 {
		timesig.RenderGregorian(ctx, canv, lines, timeSignature, x)
	}
	canv.Gend()
	canv.Gend()
	return maxY - lines[4]
}

func GetLeftIndentWithTimeSignature(key keysig.Key, timeSig timesig.TimeSignature) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return constant.LAYOUT_INDENT_LENGTH + CLEF_WIDTH + (timesig.GREGORIAN_WIDTH * len(timeSig.UniqueSign)) + (PADDING_WIDTH*(3+(len(timeSig.UniqueSign)-1)) + keySigWith)
}

func GetLeftIndent(key keysig.Key) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	offset := 0
	if key.Start && key.Prev != nil {
		offset = (len(key.Prev.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH) + PADDING_WIDTH
	}
	return constant.LAYOUT_INDENT_LENGTH + CLEF_WIDTH + (PADDING_WIDTH * 2) + keySigWith + offset
}

func GetLeftMargin(key keysig.Key) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return CLEF_WIDTH + (PADDING_WIDTH * 2) + keySigWith
}
