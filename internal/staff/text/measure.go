package text

import (
	"unicode"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/staff/lines"
)

// upper = 10 lower 5. Refrein 47. Fine 30

var defaultTextWidth = map[string]int{
	"Refrein": REFREIN_TEXT_LENGTH,
	"Fine":    FINE_TEXT_LENGTH,
}

func GetTextMarginBottom(stafflines lines.LineStaff, notes []*entity.NoteRenderer, i int) float64 {
	left := notes[i].PositionX

	maxTextWidth := 0
	leftOffset := 0
	isRightAlignment := false
	for _, t := range notes[i].MeasureText {
		if t.TextAlignment == musicxml.TextAlignmentRight {
			isRightAlignment = true
		}
		textWidth := 0
		def, ok := defaultTextWidth[t.Text]
		if ok {
			textWidth = def
		} else {
			for _, c := range t.Text {
				if unicode.IsUpper(c) {
					textWidth += OTHER_TEXT_UPPERCASE_WIDTH
				} else {
					textWidth += OTHER_TEXT_LOWERCASE_WIDTH
				}
			}
		}

		if maxTextWidth < textWidth {
			maxTextWidth = textWidth
		}
	}

	right := left + maxTextWidth
	if isRightAlignment {
		leftBound := left - maxTextWidth
		right = constant.LAYOUT_WIDTH - constant.LAYOUT_INDENT_LENGTH
		for pos := i; pos >= 0; pos-- {
			note := notes[pos]
			if note.PositionX > leftBound {
				leftOffset++
				left = note.PositionX
			}
		}

		if i-leftOffset >= 0 {
			i -= leftOffset

		}
	}

	result := 0
	for pos := i; pos < len(notes); pos++ {
		note := notes[pos]

		if note.PositionX >= right && len(note.MeasureDash) == 0 {
			break
		}

		if result < note.MarginTopFromStaff {
			result = note.MarginTopFromStaff
		}
	}
	if result <= TEXT_STAFF_TOP_LINE_GAP_WIDTH {
		return 0
	}

	return float64(result)
}
