package numbered

import (
	"context"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
)

type renderedNote struct {
	IsDotted bool
	Type     musicxml.NoteLength
}

//RenderLengthNote give the data that needed for the numbered
// TODO: add support got 8th beat and more
func RenderLengthNote(ctx context.Context, ts timesig.TimeSignature, measure int, noteLength float64) []renderedNote {

	currentTimeSig := ts.GetTimesignatureOnMeasure(ctx, measure)

	if currentTimeSig.BeatType == 4 {
		result := []renderedNote{
			renderedNote{
				Type: musicxml.NoteLengthQuarter,
			},
		}

		if noteLength == 0.5 {
			return []renderedNote{
				renderedNote{
					Type: musicxml.NoteLengthEighth,
				},
			}
		}

		for i := 1; i <= int(math.Trunc(noteLength))-1; i++ {
			result = append(result, renderedNote{IsDotted: true})
		}

		if math.Trunc(noteLength) != noteLength { // decimal dotted beat
			result = append(result, renderedNote{IsDotted: true, Type: musicxml.NoteLengthEighth})
		}

		return result

	}

	return nil

}
