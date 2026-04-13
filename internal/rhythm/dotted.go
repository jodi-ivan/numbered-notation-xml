package rhythm

import (
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/breathpause"
	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

// TODO: remove the dot positioning operation here, since it handled in the align justify.
func (ri *rhythmInteractor) AdjustMultiDottedRenderer(notes []*entity.NoteRenderer, x int, y int) (int, int) {

	xNotes := 0
	continueDot := false
	lastDotLoc := 0
	dotCount := 0

	hasDashedSlur := false
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

		} else if breathpause.IsBreathMark(n) && prev != nil && prev.IsLengthTakenFromLyric {
			x -= constant.LOWERCASE_LENGTH
		} else {
			xNotes = x
			continueDot = false
			dotCount = 0
		}
		n.PositionX = x
		n.PositionY = y
		if n.IsLengthTakenFromLyric {
			x += n.Width
		} else {
			x += int(math.Min(constant.UPPERCASE_LENGTH*2, float64(n.Width)))
			// x += 4
		}
		if prev != nil && prev.IsLengthTakenFromLyric && n.IsDotted {
			if float64(prev.Width) > float64(constant.UPPERCASE_LENGTH*dotCount) {
				diff := (prev.Width - (constant.UPPERCASE_LENGTH * dotCount))
				x = (x - diff) + constant.UPPERCASE_LENGTH
			}

		}

		for _, s := range n.Slur {
			if s.LineType == musicxml.NoteSlurLineTypeDashed {
				hasDashedSlur = true
				break
			}
		}
		// dont merge notes if it has dashed slur
		if !hasDashedSlur && (n.Tie != nil || len(n.Slur) > 0) {
			x -= n.Width
			x += constant.LOWERCASE_LENGTH * 2
			n.Width = constant.UPPERCASE_LENGTH
			hasDashedSlur = false
		}

		prev = n
		if n.IsDotted && i == len(notes)-1 && dotCount > 1 {
			x += constant.LOWERCASE_LENGTH
		}
		if n.IsNewLine {
			x = constant.LAYOUT_INDENT_LENGTH
			// y is not added up because it will handled by the staff (the function that call this function).
		}
	}

	for i, rev := range revisionX {
		note := notes[i]

		note.PositionX = rev
		notes[i] = note
	}

	return x, y
}
