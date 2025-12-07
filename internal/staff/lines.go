package staff

import (
	"context"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
)

// SplitLines split the measure in the lines manner
func (si *staffInteractor) SplitLines(ctx context.Context, part musicxml.Part) [][]musicxml.Measure {
	result := [][]musicxml.Measure{}
	currentLine := []musicxml.Measure{}
	isLastMeasure := false
	for i, measure := range part.Measures {

		if measure.Print != nil && measure.Print.NewSystem == musicxml.PrintNewSystemTypeYes {
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
