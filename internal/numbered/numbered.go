package numbered

import (
	"context"
	"math"

	"github.com/jodi-ivan/numbered-notation-xml/internal/entity"
	"github.com/jodi-ivan/numbered-notation-xml/internal/musicxml"
	"github.com/jodi-ivan/numbered-notation-xml/internal/timesig"
	"github.com/jodi-ivan/numbered-notation-xml/utils/canvas"
)

type NoteLength struct {
	IsDotted bool
	Type     musicxml.NoteLength
}

type Numbered interface {
	GetLengthNote(ctx context.Context, ts timesig.TimeSignature, measure int, noteLength float64) []NoteLength
	RenderOctave(ctx context.Context, canv canvas.Canvas, notes []*entity.NoteRenderer)
}

type numberedInteractor struct{}

func (ni *numberedInteractor) GetLengthNote(ctx context.Context, ts timesig.TimeSignature, measure int, noteLength float64) []NoteLength {
	return RenderLengthNote(ctx, ts, measure, noteLength)
}

func (ni *numberedInteractor) RenderOctave(ctx context.Context, canv canvas.Canvas, notes []*entity.NoteRenderer) {
	RenderOctave(ctx, canv, notes)
}

func New() Numbered {
	return &numberedInteractor{}
}

// RenderLengthNote give the data that needed for the numbered
func RenderLengthNote(ctx context.Context, ts timesig.TimeSignature, measure int, noteLength float64) []NoteLength {

	currentTimeSig := ts.GetTimesignatureOnMeasure(ctx, measure)

	if currentTimeSig.BeatType == 4 {
		result := []NoteLength{
			NoteLength{
				Type: musicxml.NoteLengthQuarter,
			},
		}

		if noteLength == 0.75 {
			return []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthEighth,
				},
				NoteLength{IsDotted: true, Type: musicxml.NoteLength16th},
			}
		}

		if noteLength == 0.5 {
			return []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthEighth,
				},
			}
		}

		if noteLength == 0.25 {
			return []NoteLength{
				NoteLength{
					Type: musicxml.NoteLength16th,
				},
			}
		}

		for i := 1; i <= int(math.Trunc(noteLength))-1; i++ {
			result = append(result, NoteLength{IsDotted: true, Type: musicxml.NoteLengthQuarter})
		}

		if math.Trunc(noteLength) != noteLength { // decimal dotted beat
			result = append(result, NoteLength{IsDotted: true, Type: musicxml.NoteLengthEighth})
		}

		return result

	}

	if currentTimeSig.BeatType == 8 {
		result := []NoteLength{}

		if noteLength == 1 {
			return []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthEighth,
				},
			}
		}

		if noteLength == 0.75 { // 3
			return []NoteLength{
				NoteLength{
					Type: musicxml.NoteLengthEighth,
				},
				NoteLength{
					IsDotted: true,
				},
			}
		}

		if noteLength == 0.5 {
			return []NoteLength{
				NoteLength{
					Type: musicxml.NoteLength16th,
				},
			}
		}

		if noteLength == 0.25 {
			return []NoteLength{
				NoteLength{
					Type: musicxml.NoteLength32nd,
				},
			}
		}

		if noteLength == 0.125 {
			return []NoteLength{
				NoteLength{
					Type: musicxml.NoteLength64th,
				},
			}
		}

		if noteLength == 0.0625 {
			return []NoteLength{
				NoteLength{
					Type: musicxml.NoteLength128th,
				},
			}
		}

		result = append(result, NoteLength{
			Type: musicxml.NoteLengthEighth,
		})
		for i := 1; i <= int(math.Trunc(noteLength))-1; i++ {
			result = append(result, NoteLength{IsDotted: true, Type: musicxml.NoteLengthEighth})
		}

		if math.Trunc(noteLength) != noteLength { // decimal dotted beat
			result = append(result, NoteLength{IsDotted: true, Type: musicxml.NoteLengthEighth})
		}

		return result
	}

	return nil

}
