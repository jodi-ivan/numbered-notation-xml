package gregorian

import (
	"context"
	"fmt"
	"unicode"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
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
	x2 := constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH + 4
	canv.Group(`class="staff-line"`)
	for i := 0; i <= 4; i++ {
		lines[i] = y - 70
		canv.Line(constant.LAYOUT_INDENT_LENGTH, y-70, x2, y-70, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.8")
		y += 8
	}
	canv.Line(constant.LAYOUT_INDENT_LENGTH, initialY, constant.LAYOUT_INDENT_LENGTH, y-70-8, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1")
	canv.Gend()

	canv.Group(`class="notes"`, `style="font-size:2em"`)
	for _, note := range notes {
		if note.IsAdditional {
			continue
		}

		if note.AbsoluteNote != "" {
			bean := `&#xF064;`
			if note.NoteLength == musicxml.NoteLengthHalf {
				bean = `&#xF063;`
			}
			canv.TextUnescaped(float64(note.PositionX), GetYpos(lines, 8, note.AbsoluteOctave, rune(note.AbsoluteNote[0])),
				bean,
				fmt.Sprintf(`pitch="%s"`, note.AbsoluteNote), fmt.Sprintf(`octave="%d"`, note.AbsoluteOctave))
		}

	}
	canv.Gend()

	x := float64(constant.LAYOUT_INDENT_LENGTH)

	// clef
	key := keySignature.GetKeyOnMeasure(ctx, 1)
	accidentalSet := key.GetAccidentals()

	canv.Group(`class="keysig"`, `style="font-size:1.8em"`)
	for x, acc := range accidentalSet {
		accidental := `&#xF02B;`
		width := 8.0
		if key.Fifth < 0 {
			accidental = `&#xF02D;`
		}
		canv.TextUnescaped(float64(constant.LAYOUT_INDENT_LENGTH+35)+(width*float64(x)),
			GetYPosKeySig(lines, 8, acc, key.Fifth < 0),
			accidental)
	}
	canv.Gend()

	canv.Group(`class="clef"`, `style="font-size:2em"`)
	canv.TextUnescaped(constant.LAYOUT_INDENT_LENGTH+8, float64(initialY+15), `&#xF026;`)
	canv.Gend()

	x += 35 + float64(len(accidentalSet)*8)

	if staffPos == 0 {
		canv.Group(`class="timesig"`, `style="font-size:2em"`)
		canv.TextUnescaped(x+10, float64(lines[0]+lines[2])/2, `&#xF034;`)
		canv.TextUnescaped(x+10, float64(lines[2]+lines[4])/2, `&#xF034;`)
		canv.Gend()
		x += 10

	}

	canv.Gend()
	return 0
}

func GetLeftIndentWithTimeSignature(key keysig.Key) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return constant.LAYOUT_INDENT_LENGTH + CLEF_WIDTH + timesig.GREGORIAN_WIDTH + (PADDING_WIDTH * 3) + keySigWith
}

func GetLeftIndent(key keysig.Key) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return constant.LAYOUT_INDENT_LENGTH + CLEF_WIDTH + (PADDING_WIDTH * 2) + keySigWith
}

func GetLeftMargin(key keysig.Key) int {
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return CLEF_WIDTH + (PADDING_WIDTH * 2) + keySigWith
}
