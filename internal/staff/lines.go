package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/gregorian"
	"github.com/jodi-ivan/numbered-notation-xml/internal/keysig"
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

func ProcessPreviousLines(prevNotes []*entity.NoteRenderer, ks keysig.KeySignature, yPos int) ([][]*entity.NoteRenderer, StaffInfo) {
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
			key := ks.GetKeyOnMeasure(context.Background(), prevNotes[0].MeasureNumber) // TEST: need to be tested
			staffInfo.MarginLeft = gregorian.GetLeftIndent(key)
			prevNotes[pos].IsNewLine = false
			staffInfo.ForceNewLine = true
		}
		staffInfo.NextLineRenderer = prevNotes[pos+offset:]

		staffInfo.Multiline = true
	} else {
		result = append(result, prevNotes)

		key := ks.GetKeyOnMeasure(context.Background(), prevNotes[0].MeasureNumber)
		staffInfo.MarginLeft = gregorian.GetLeftIndent(key)
	}

	return result, staffInfo
}

func PrepareNextLines(staffInfo StaffInfo, ks keysig.KeySignature, notes []*entity.NoteRenderer, rightBarline *entity.NoteRenderer) StaffInfo {
	proceed := false
	maxTotalLyric := 1

	key := ks.GetKeyOnMeasure(context.Background(), notes[0].MeasureNumber)

	indent := gregorian.GetLeftIndent(key)
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
			note.PositionX = indent
		}
		staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, note)
	}
	if staffInfo.MarginBottom < (maxTotalLyric-1)*25 {
		staffInfo.MarginBottom = (maxTotalLyric - 1) * 25
	}

	if len(staffInfo.NextLineRenderer) == 0 && rightBarline.Barline.BarStyle == musicxml.BarLineStyleNone {
		return StaffInfo{
			NextLineRenderer: []*entity.NoteRenderer{},
			MarginLeft:       indent,
		}
	}
	staffInfo.NextLineRenderer = append(staffInfo.NextLineRenderer, rightBarline)
	return staffInfo

}
