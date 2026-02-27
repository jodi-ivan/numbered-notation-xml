package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/constant"
	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

// SplitLines split the measure in the lines manner
func (si *staffInteractor) SplitLines(ctx context.Context, part musicxml.Part) [][]musicxml.Measure {
	result := [][]musicxml.Measure{}
	currentLine := []musicxml.Measure{}
	isLastMeasure := false
	for i, measure := range part.Measures {

		if measure.Print != nil && (measure.Print.NewSystem == musicxml.PrintNewSystemTypeYes || measure.Print.NewPage == musicxml.PrintNewSystemTypeYes) {
			isLastMeasure = i == (len(part.Measures) - 1)
			finishLine := make([]musicxml.Measure, len(currentLine))
			copy(finishLine, currentLine)

			result = append(result, finishLine)

			currentLine = []musicxml.Measure{}
		}
		currentLine = append(currentLine, measure)

	}

	if isLastMeasure {
		result = append(result, currentLine)
		return append(result, []musicxml.Measure{})
	}

	return append(result, currentLine)
}

func ProcessPreviousLines(prevNotes []*entity.NoteRenderer, yPos int) ([][]*entity.NoteRenderer, StaffInfo) {
	result := [][]*entity.NoteRenderer{}
	staffInfo := StaffInfo{}
	pos := -1

	// last line with no staff measure remaining
	for i, note := range prevNotes {
		note.PositionY = yPos
		if note.IsNewLine {
			pos = i
			break
		}
	}

	if pos != -1 {
		result = append(result, prevNotes[:pos+1])
		staffInfo.NextLineRenderer = prevNotes[pos+1:]
		staffInfo.MarginLeft = constant.LAYOUT_INDENT_LENGTH
		staffInfo.Multiline = true
		// TODO: assign a proper margin botton (multiline lyric)
		// staffInfo.MarginBottom = 80
	} else {
		result = append(result, prevNotes)
		staffInfo.MarginLeft = constant.LAYOUT_INDENT_LENGTH
	}

	return result, staffInfo
}

func PrepareNextLines(staffInfo StaffInfo, notes []*entity.NoteRenderer, rightBarline *entity.NoteRenderer) StaffInfo {
	proceed := false
	for _, note := range notes {
		if !proceed {
			if note.IsNewLine {
				proceed = true
			}
			continue
		}
		if len(staffInfo.NextLineRenderer) == 0 {
			note.PositionX = constant.LAYOUT_INDENT_LENGTH
		}
		staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, note)
	}
	staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, rightBarline)
	return staffInfo

}
