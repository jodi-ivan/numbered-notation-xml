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

func RenderStaffLine(ctx context.Context, staffPos, y int, canv canvas.Canvas, notes []*entity.NoteRenderer, keySignature keysig.KeySignature, timeSignature timesig.TimeSignature) int {
	initialY := y - 70
	lines := [5]int{}
	canv.Group(`class="gregorian"`, "style='font-family:mozart11'")
	x2 := constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH + 8
	canv.Group(`class="staff-line"`)
	for i := 0; i <= 4; i++ {
		lines[i] = y - 70
		canv.Line(constant.LAYOUT_INDENT_LENGTH, y-70, x2, y-70, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.8")
		y += STAFF_SPACE_WIDTH
	}
	canv.Line(constant.LAYOUT_INDENT_LENGTH, initialY, constant.LAYOUT_INDENT_LENGTH, y-70-STAFF_SPACE_WIDTH, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1")
	canv.Gend()

	maxY := lines[4]

	groupBeam := [][]CoordinateWithNoteLength{{}}

	canv.Group(`class="notes"`, `style="font-size:2em"`)
	currentMeasure := 0

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
			xPos := note.PositionX
			if note.PositionX-notes[i-1].PositionX <= numbered.MIN_DISTANCE_BREATH {
				xPos += (numbered.AVERAGE_CHARACTER_WIDTH + constant.LOWERCASE_LENGTH) / 3
			}
			canv.TextUnescaped(float64(xPos), float64(lines[0])-8, "&#xF0E2;", `style="font-size:1.3em"`)
			continue
		}
		if note.IsRest {
			canv.TextUnescaped(float64(note.PositionX), float64(lines[2]), restHex[note.NoteLength])
			continue
		}

		if note.Barline != nil {
			barline.RenderGregorian(canv, note.Barline, i == len(notes)-1, lines, entity.NewCoordinate(float64(note.PositionX), float64(lines[4])))
			continue
		}

		if note.AbsoluteNote == "" {
			continue
		}

		yPos := GetYpos(lines, STAFF_SPACE_WIDTH, note.AbsoluteOctave, rune(note.AbsoluteNote[0]))
		if maxY < int(yPos) {
			maxY = int(yPos)
		}

		RenderLedgerLine(canv, entity.NewCoordinate(float64(note.PositionX), yPos), lines)

		canv.TextUnescaped(float64(note.PositionX), yPos,
			beanNoteHex[note.NoteLength],
			fmt.Sprintf(`pitch="%s"`, note.AbsoluteNote), fmt.Sprintf(`octave="%d"`, note.AbsoluteOctave))

		dottedHalf := i < len(notes)-1 && notes[i+1].IsDotted && len(notes[i+1].Beam) >= 1
		singleNoteValue := timeSignature.GetNoteLength(ctx, note.MeasureNumber, musicxml.Note{Type: note.NoteLength})
		dottedBeat := note.NoteValue > singleNoteValue && note.Tie == nil
		if dottedHalf || dottedBeat {
			dotPos := yPos

			if (int(yPos)-initialY)%STAFF_SPACE_WIDTH == 0 {
				dotPos -= 4
			}
			canv.TextUnescaped(float64(note.PositionX+15), dotPos, "&#xF060;")
		}

		if note.NoteLength == musicxml.NoteLengthWhole {
			continue
		}

		if len(note.Beam) == 0 {
			renderMap[cmp.Compare(yPos, float64(lines[2]))](canv, lines, CoordinateWithNoteLength{Coordinate: entity.NewCoordinate(float64(note.PositionX), yPos), NoteLength: note.NoteLength})
			continue
		}

		groupBeam[len(groupBeam)-1] = append(groupBeam[len(groupBeam)-1], CoordinateWithNoteLength{
			Coordinate: entity.NewCoordinate(float64(note.PositionX), yPos),
			NoteLength: note.NoteLength,
		})
		if len(note.Beam) >= 1 && note.Beam[1].Type == musicxml.NoteBeamTypeEnd {
			groupBeam = append(groupBeam, []CoordinateWithNoteLength{})
		}

	}

	for _, gr := range groupBeam {
		if len(gr) == 0 {
			continue
		}
		RenderGroupBeam(canv, gr, lines)
	}
	canv.Gend()
	canv.Gend()

	x := float64(constant.LAYOUT_INDENT_LENGTH)

	canv.Group(`class="staff-markings"`)
	// clef
	key := keySignature.GetKeyOnMeasure(ctx, 1)
	accidentalSet := key.GetAccidentals()

	canv.Group(`class="keysig"`, `style="font-size:1.8em"`)
	for x, acc := range accidentalSet {
		accidental := `&#xF02B;` // sharp
		width := 8.0
		if key.Fifth < 0 {
			accidental = `&#xF02D;` // flat
		}
		canv.TextUnescaped(float64(constant.LAYOUT_INDENT_LENGTH+35)+(width*float64(x)),
			GetYPosKeySig(lines, 8, acc, key.Fifth < 0),
			accidental)
	}
	canv.Gend()

	canv.Group(`class="clef"`, `style="font-size:2em"`)
	canv.TextUnescaped(constant.LAYOUT_INDENT_LENGTH+5, float64(initialY+15), `&#xF026;`)
	canv.Gend()

	x += 35 + float64(len(accidentalSet)*ACCIDENTAL_KEY_SIGNATURE_WIDTH) + PADDING_WIDTH

	if staffPos == 0 {
		timesig.RenderGregorian(ctx, canv, lines, timeSignature, x)
	}
	canv.Gend()
	canv.Gend()
	return maxY
}

func GetLeftIndentWithTimeSignature(key keysig.Key, timeSig timesig.TimeSignature) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return constant.LAYOUT_INDENT_LENGTH + CLEF_WIDTH + (timesig.GREGORIAN_WIDTH * len(timeSig.UniqueSign)) + (PADDING_WIDTH*(3+(len(timeSig.UniqueSign)-1)) + keySigWith)
}

func GetLeftIndent(key keysig.Key) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return constant.LAYOUT_INDENT_LENGTH + CLEF_WIDTH + (PADDING_WIDTH * 2) + keySigWith
}

func GetLeftMargin(key keysig.Key) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return CLEF_WIDTH + (PADDING_WIDTH * 2) + keySigWith
}
