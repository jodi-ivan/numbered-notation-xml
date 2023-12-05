package rhythm

import (
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
)

func SetAndAdjustMultiDottedRenderer(notes []*entity.NoteRenderer, x int, y int) (int, int, bool, map[int]int) {
	xNotes := 0
	continueDot := false
	lastDotLoc := 0
	dotCount := 0

	multiline := false

	var prev *entity.NoteRenderer

	revisionX := map[int]int{}
	for i, n := range notes {
		if n.IsDotted {
			dotCount++
			if continueDot {
				revisionX[i] = lastDotLoc + constant.UPPERCASE_LENGTH
				lastDotLoc = lastDotLoc + constant.UPPERCASE_LENGTH
			} else {
				revisionX[i] = xNotes + constant.UPPERCASE_LENGTH
				lastDotLoc = xNotes + constant.UPPERCASE_LENGTH
			}
			continueDot = true
		} else if n.Articulation != nil && n.Articulation.BreathMark != nil {
			if prev != nil && prev.IsLengthTakenFromLyric {
				x -= prev.Width - constant.LOWERCASE_LENGTH
			}
		} else {
			if continueDot {
				x += constant.LOWERCASE_LENGTH
			}
			xNotes = x
			continueDot = false
			dotCount = 0
		}

		n.PositionX = x
		n.PositionY = y
		x += n.Width
		if prev != nil && prev.IsLengthTakenFromLyric && n.IsDotted {
			x = x - n.Width
		}
		if n.IsNewLine {
			x = constant.LAYOUT_INDENT_LENGTH
			multiline = multiline || true
		}
		n.IndexPosition = i
		prev = n
		if n.IsDotted && i == len(notes)-1 && dotCount > 1 {
			x += constant.LOWERCASE_LENGTH
		}
	}

	return x, y, multiline, revisionX
}
