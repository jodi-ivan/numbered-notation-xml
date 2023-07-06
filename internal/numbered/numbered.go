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

// RenderLengthNote give the data that needed for the numbered
func RenderLengthNote(ctx context.Context, ts timesig.TimeSignature, measure int, noteLength float64) []renderedNote {

	currentTimeSig := ts.GetTimesignatureOnMeasure(ctx, measure)

	if currentTimeSig.BeatType == 4 {
		result := []renderedNote{
			renderedNote{
				Type: musicxml.NoteLengthQuarter,
			},
		}

		if noteLength == 0.75 {
			return []renderedNote{
				renderedNote{
					Type: musicxml.NoteLengthEighth,
				},
				renderedNote{IsDotted: true, Type: musicxml.NoteLength16th},
			}
		}

		if noteLength == 0.5 {
			return []renderedNote{
				renderedNote{
					Type: musicxml.NoteLengthEighth,
				},
			}
		}

		if noteLength == 0.25 {
			return []renderedNote{
				renderedNote{
					Type: musicxml.NoteLength16th,
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

	if currentTimeSig.BeatType == 8 {
		result := []renderedNote{}

		if noteLength == 1 {
			return []renderedNote{
				renderedNote{
					Type: musicxml.NoteLengthEighth,
				},
			}
		}

		if noteLength == 0.75 { // 3
			return []renderedNote{
				renderedNote{
					Type: musicxml.NoteLengthEighth,
				},
				renderedNote{
					IsDotted: true,
				},
			}
		}

		if noteLength == 0.5 {
			return []renderedNote{
				renderedNote{
					Type: musicxml.NoteLength16th,
				},
			}
		}

		if noteLength == 0.25 {
			return []renderedNote{
				renderedNote{
					Type: musicxml.NoteLength32nd,
				},
			}
		}

		if noteLength == 0.125 {
			return []renderedNote{
				renderedNote{
					Type: musicxml.NoteLength64th,
				},
			}
		}

		if noteLength == 0.0625 {
			return []renderedNote{
				renderedNote{
					Type: musicxml.NoteLength128th,
				},
			}
		}

		result = append(result, renderedNote{
			Type: musicxml.NoteLengthEighth,
		})
		for i := 1; i <= int(math.Trunc(noteLength))-1; i++ {
			result = append(result, renderedNote{IsDotted: true, Type: musicxml.NoteLengthEighth})
		}

		if math.Trunc(noteLength) != noteLength { // decimal dotted beat
			result = append(result, renderedNote{IsDotted: true, Type: musicxml.NoteLengthEighth})
		}

		return result
	}

	return nil

}
