package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

// SplitLines split the measure in the lines manner
func (si *staffInteractor) SplitLines(ctx context.Context, part musicxml.Part) [][]musicxml.Measure {
	result := [][]musicxml.Measure{}
	currentLine := []musicxml.Measure{}
	for _, measure := range part.Measures {

		if measure.Print != nil && measure.Print.NewSystem == musicxml.PrintNewSystemTypeYes {
			finishLine := make([]musicxml.Measure, len(currentLine))
			copy(finishLine, currentLine)

			result = append(result, finishLine)

			currentLine = []musicxml.Measure{}
		}
		currentLine = append(currentLine, measure)

	}

	return append(result, currentLine)
}
