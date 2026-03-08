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
	maxTotalLyric := 1

	// last line with no staff measure remaining
	for i, note := range prevNotes {
		note.PositionY = yPos
		if maxTotalLyric < len(note.Lyric) {
			maxTotalLyric = len(note.Lyric)
		}
		if note.IsNewLine {
			pos = i
			break
		}
	}

	staffInfo.MarginBottom = (maxTotalLyric - 1) * 25
	if pos != -1 {
		result = append(result, prevNotes[:pos+1])
		offset := 1

		if pos+1 < len(prevNotes) && prevNotes[pos+1].Barline != nil {
			// kj-139, a whole new line for the next renderer and the current line transffered to the next-next line
			barlineNote := prevNotes[pos+1]
			barlineNote.PositionX = prevNotes[pos].PositionX + prevNotes[pos].Width
			result[0] = append(result[0], barlineNote)

			offset = 2
			staffInfo.MarginLeft = constant.LAYOUT_INDENT_LENGTH
			prevNotes[pos].IsNewLine = false
			staffInfo.ForceNewLine = true
		}
		staffInfo.NextLineRenderer = prevNotes[pos+offset:]

		staffInfo.Multiline = true
	} else {
		result = append(result, prevNotes)
		staffInfo.MarginLeft = constant.LAYOUT_INDENT_LENGTH
	}

	return result, staffInfo
}

func PrepareNextLines(staffInfo StaffInfo, notes []*entity.NoteRenderer, rightBarline *entity.NoteRenderer) StaffInfo {
	proceed := false
	maxTotalLyric := 1

	for _, note := range notes {
		if !proceed {
			if note.IsNewLine {
				proceed = true
			}
			continue
		}

		if maxTotalLyric < len(note.Lyric) {
			maxTotalLyric = len(note.Lyric)
		}

		if len(staffInfo.NextLineRenderer) == 0 {
			note.PositionX = constant.LAYOUT_INDENT_LENGTH
		}
		staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, note)
	}
	if staffInfo.MarginBottom < (maxTotalLyric-1)*25 {
		staffInfo.MarginBottom = (maxTotalLyric - 1) * 25
	}
	staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, rightBarline)
	return staffInfo

}
