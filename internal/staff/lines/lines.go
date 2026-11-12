package lines

import (
	"cmp"
	"context"
	"unicode"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
	"github.com/jodi-ivan/numbered-notation-xml/internal/layout"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type LineStaff interface {
	GetLines() [5]int
	GetMiddleLine() int
	GetTopLine() int
	GetBottomLine() int
	GetStemDirection(pitch rune, octave int) int
	GetStemDirectionCompare(yPos float64) int
	GetYPosKeySig(pitch string, isFlat bool) float64
	Render(ctx context.Context, canv canvas.Canvas, y int, measureNo int, inclTimesig bool)
	GetYPos(pitch rune, octave int) float64
	GetLeftIndent(ctx context.Context, measure ...int) int
	SetMarginRight(mr int)
	GetMarginRight() int
	GetLeftIndentWithTimeSignature(ctx context.Context) int
	GetKeysig() keysig.KeySignature
	GetTimeSig() timesig.TimeSignature
}

type lineStaff struct {
	Lines        [5]int
	Keysig       keysig.KeySignature
	TimeSig      timesig.TimeSignature
	MarginRight  int
	LeftIndent   int
	MeasureStart int
}

func NewMiddleNonFirstLineStaff(ks keysig.KeySignature) LineStaff {
	return &lineStaff{
		Keysig: ks,
	}
}

func NewLineStaff(ts timesig.TimeSignature, ks keysig.KeySignature) LineStaff {
	return &lineStaff{
		Keysig:  ks,
		TimeSig: ts,
	}
}

func NewLineStaffWithLines(ts timesig.TimeSignature, ks keysig.KeySignature, y int) LineStaff {

	result := &lineStaff{
		Keysig:  ks,
		TimeSig: ts,
	}
	for i := 0; i <= 4; i++ {
		result.Lines[i] = y
		y += STAFF_SPACE_WIDTH
	}

	return result
}
func (ls *lineStaff) GetLines() [5]int {
	return ls.Lines
}

func (ls *lineStaff) GetMiddleLine() int {
	return ls.Lines[2]
}

func (ls *lineStaff) GetTopLine() int {
	return ls.Lines[0]
}

func (ls *lineStaff) GetBottomLine() int {
	return ls.Lines[4]
}

func (ls *lineStaff) GetStemDirection(pitch rune, octave int) int {
	yPos := ls.GetYPos(pitch, octave)
	direction := cmp.Compare(yPos, float64(ls.GetMiddleLine()))
	if direction == 0 {
		direction = -1
	}

	return direction
}

func (ls *lineStaff) GetStemDirectionCompare(yPos float64) int {
	direction := cmp.Compare(yPos, float64(ls.GetMiddleLine()))
	if direction == 0 {
		direction = -1
	}

	return direction

}

func (ls *lineStaff) GetYPosKeySig(pitch string, isFlat bool) float64 {
	if isFlat {
		// Flat order: B E A D G C F
		pos := map[string]float64{
			"B": float64(ls.Lines[2]),
			"E": float64(ls.Lines[0]) + (constant.SPACE_LENGTH / 2),
			"A": float64(ls.Lines[2]) + (constant.SPACE_LENGTH / 2),
			"D": float64(ls.Lines[1]),
			"G": float64(ls.Lines[3]),
			"C": float64(ls.Lines[1]) + (constant.SPACE_LENGTH / 2),
			"F": float64(ls.Lines[3]) - (constant.SPACE_LENGTH / 2),
		}
		return pos[pitch]
	}
	// Sharp order: F C G D A E B
	pos := map[string]float64{
		"F": float64(ls.Lines[0]),
		"C": float64(ls.Lines[2]) - (constant.SPACE_LENGTH / 2),
		"G": float64(ls.Lines[1]) - (constant.SPACE_LENGTH / 2),
		"D": float64(ls.Lines[2]),
		"A": float64(ls.Lines[3]) + (constant.SPACE_LENGTH / 2),
		"E": float64(ls.Lines[0]) + (constant.SPACE_LENGTH / 2),
		"B": float64(ls.Lines[2]) - (constant.SPACE_LENGTH),
	}
	return pos[pitch]
}

func (ls *lineStaff) Render(ctx context.Context, canv canvas.Canvas, y int, measureNo int, inclTimesig bool) {
	if ls.MarginRight == 0 {
		ls.MarginRight = constant.LAYOUT_WIDTH
	}
	leftMargin, rightMargin := layout.GetMargin(ctx)
	ls.MarginRight -= rightMargin

	ls.MeasureStart = measureNo
	canv.Group(`class="staff-line"`)
	for i := 0; i <= 4; i++ {
		ls.Lines[i] = y
		canv.Line(leftMargin, y, ls.MarginRight, y, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:0.8")
		y += STAFF_SPACE_WIDTH
	}
	canv.Line(leftMargin, ls.Lines[0], leftMargin, y-STAFF_SPACE_WIDTH, "fill:none;stroke:#000000;stroke-linecap:round;stroke-width:1")
	canv.Gend()

	x := float64(leftMargin)
	initialY := ls.Lines[0]
	canv.Group(`class="staff-markings"`)
	// clef
	key := ls.Keysig.GetKeyOnMeasure(context.Background(), measureNo)
	accidentalSet := key.GetAccidentals()

	canv.Group(`class="clef"`, `style="font-size:28px"`)
	canv.TextUnescaped(float64(leftMargin)+5, float64(initialY+15), TREBLE_CLEF_HEX)
	canv.Gend()

	canv.Group(`class="keysig"`, `style="font-size:28px"`)
	offset := 0

	// key signature changes
	if key.Start && key.Prev != nil && measureNo != 1 {
		naturalSet := key.Prev.GetAccidentals()

		for x, acc := range naturalSet {
			accidental := accidentalHex[musicxml.NoteAccidentalNatural]
			width := ACCIDENTAL_KEY_SIGNATURE_WIDTH

			canv.TextUnescaped(float64(leftMargin+CLEF_WIDTH)+float64(width*x),
				ls.GetYPosKeySig(acc, key.Prev.Fifth < 0),
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
		canv.TextUnescaped(float64(leftMargin+CLEF_WIDTH+offset)+float64(width*x),
			ls.GetYPosKeySig(acc, key.Fifth < 0),
			accidental)
	}
	canv.Gend()

	x += CLEF_WIDTH + (float64(len(accidentalSet)) * ACCIDENTAL_KEY_SIGNATURE_WIDTH) + PADDING_WIDTH + float64(offset)

	// if staffPos == 0 && len(timeSignature.Signatures) > 0 {
	if inclTimesig && len(ls.TimeSig.Signatures) > 0 {
		timesig.RenderGregorian(context.Background(), canv, ls.Lines, ls.TimeSig, x)
	}
	canv.Gend()
}

func (ls *lineStaff) GetYPos(pitch rune, octave int) float64 {
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

	return float64(ls.Lines[0]) + float64(stepsBelow)*(float64(STAFF_SPACE_WIDTH)/2)
}

func (ls *lineStaff) GetLeftIndent(ctx context.Context, measure ...int) int {
	measureNum := ls.MeasureStart
	if len(measure) > 0 {
		measureNum = measure[0]
	}

	key := ls.Keysig.GetKeyOnMeasure(context.Background(), measureNum)
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	offset := 0
	if key.Start && key.Prev != nil {
		offset = (len(key.Prev.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH) + PADDING_WIDTH
	}
	leftMargin, _ := layout.GetMargin(ctx)
	return leftMargin + CLEF_WIDTH + (PADDING_WIDTH * 2) + keySigWith + offset
}

func (ls *lineStaff) SetMarginRight(mr int) {
	ls.MarginRight = mr
}

func (ls *lineStaff) GetMarginRight() int {
	return ls.MarginRight
}

func (ls *lineStaff) GetLeftIndentWithTimeSignature(ctx context.Context) int {
	leftMargin, _ := layout.GetMargin(ctx)
	key := ls.Keysig.GetKeyOnMeasure(context.Background(), 1)
	keySigWith := len(key.GetAccidentals()) * ACCIDENTAL_KEY_SIGNATURE_WIDTH
	return leftMargin + CLEF_WIDTH + (timesig.GREGORIAN_WIDTH * len(ls.TimeSig.UniqueSign)) + (PADDING_WIDTH*(3+(len(ls.TimeSig.UniqueSign)-1)) + keySigWith)
}

func (ls *lineStaff) GetKeysig() keysig.KeySignature {
	return ls.Keysig
}

func (ls *lineStaff) GetTimeSig() timesig.TimeSignature {
	return ls.TimeSig
}
